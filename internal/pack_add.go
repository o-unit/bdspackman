package internal

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const maxArchiveFiles = 10000

// AddPack imports packs from a directory or a .mcpack, .mcaddon, or .zip
// archive into the selected Server or World location. It either installs every
// discovered pack or removes any packs installed during a failed operation.
func AddPack(cfg Config, inputPath string, toWorld bool) ([]Pack, error) {
	inputPath = strings.TrimSpace(inputPath)
	if inputPath == "" {
		return nil, errors.New("pack path is empty")
	}

	inputInfo, err := os.Stat(inputPath)
	if err != nil {
		return nil, fmt.Errorf("stat pack path %s: %w", inputPath, err)
	}

	workDir, err := os.MkdirTemp("", "bdspackman-add-*")
	if err != nil {
		return nil, fmt.Errorf("create add-pack work directory: %w", err)
	}
	defer os.RemoveAll(workDir)

	root := inputPath
	if !inputInfo.IsDir() {
		root = filepath.Join(workDir, "input")
		if err := extractArchive(inputPath, root); err != nil {
			return nil, err
		}
	}
	if err := extractNestedArchives(root); err != nil {
		return nil, err
	}

	candidates, err := discoverAddPacks(root, cfg.Language)
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no pack manifest found in %s", inputPath)
	}
	if err := validateAddCandidates(cfg, candidates, toWorld); err != nil {
		return nil, err
	}

	installed := make([]Pack, 0, len(candidates))
	for _, candidate := range candidates {
		destination := filepath.Join(packDirectoryForLocation(cfg, candidate.Type, candidate.Location), candidate.FolderName)
		if err := CopyDir(candidate.Path, destination, FileOperationOptions{}); err != nil {
			for _, pack := range installed {
				_ = RemovePath(pack.Path)
			}
			return nil, fmt.Errorf("install pack %s: %w", candidate.FolderName, err)
		}
		candidate.Path = destination
		candidate.ManifestPath = filepath.Join(destination, "manifest.json")
		installed = append(installed, candidate)
	}

	return installed, nil
}

func extractNestedArchives(root string) error {
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() || !isPackArchive(path) {
			return walkErr
		}
		destination := path + ".contents"
		if err := extractArchive(path, destination); err != nil {
			return err
		}
		return nil
	})
}

func extractArchive(path, destination string) error {
	if !isPackArchive(path) {
		return fmt.Errorf("unsupported pack file: %s", path)
	}
	reader, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("open archive %s: %w", path, err)
	}
	defer reader.Close()
	if len(reader.File) > maxArchiveFiles {
		return fmt.Errorf("archive contains too many files: %s", path)
	}
	if err := EnsureDir(destination); err != nil {
		return err
	}
	for _, file := range reader.File {
		target, err := safeArchivePath(destination, file.Name)
		if err != nil {
			return fmt.Errorf("extract archive %s: %w", path, err)
		}
		if file.FileInfo().IsDir() {
			if err := EnsureDir(target); err != nil {
				return err
			}
			continue
		}
		if !file.Mode().IsRegular() {
			continue
		}
		if err := EnsureDir(filepath.Dir(target)); err != nil {
			return err
		}
		in, err := file.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_EXCL, file.Mode().Perm())
		if err != nil {
			in.Close()
			return err
		}
		_, copyErr := io.Copy(out, in)
		closeErr := out.Close()
		in.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
	}
	return nil
}

func safeArchivePath(root, name string) (string, error) {
	target := filepath.Join(root, filepath.FromSlash(name))
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	if targetAbs != rootAbs && !isSubPath(targetAbs, rootAbs) {
		return "", fmt.Errorf("archive entry escapes destination: %s", name)
	}
	return targetAbs, nil
}

func isPackArchive(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".zip", ".mcpack", ".mcaddon":
		return true
	}
	return false
}

func discoverAddPacks(root, language string) ([]Pack, error) {
	var packs []Pack
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || entry.Name() != "manifest.json" {
			return nil
		}
		manifest, err := LoadManifest(path)
		if err != nil {
			return fmt.Errorf("read manifest %s: %w", path, err)
		}
		packType, err := packTypeFromManifest(manifest)
		if err != nil {
			return fmt.Errorf("manifest %s: %w", path, err)
		}
		folder := filepath.Base(filepath.Dir(path))
		if _, err := validatePackDirectoryName(folder); err != nil {
			return err
		}
		packs = append(packs, Pack{Name: manifest.Header.Name, DisplayName: manifest.Header.Name, FolderName: folder, UUID: manifest.Header.UUID, Version: manifest.Header.Version, Type: packType, Path: filepath.Dir(path), ManifestPath: path, Status: StatusOff})
		return nil
	})
	return packs, err
}

func packTypeFromManifest(manifest *Manifest) (PackType, error) {
	if manifest.Header.UUID == "" || len(manifest.Header.Version) == 0 {
		return PackTypeBehavior, errors.New("header UUID and version are required")
	}
	for _, module := range manifest.Modules {
		switch module.Type {
		case "data":
			return PackTypeBehavior, nil
		case "resources":
			return PackTypeResource, nil
		}
	}
	return PackTypeBehavior, errors.New("does not contain a data or resources module")
}

func validateAddCandidates(cfg Config, candidates []Pack, toWorld bool) error {
	existing, err := ScanPacks(cfg)
	if err != nil {
		return fmt.Errorf("scan installed packs: %w", err)
	}
	seenUUID := make(map[PackType]map[string]struct{})
	seenFolder := make(map[PackType]map[string]struct{})
	for _, p := range existing {
		if seenUUID[p.Type] == nil {
			seenUUID[p.Type] = map[string]struct{}{}
			seenFolder[p.Type] = map[string]struct{}{}
		}
		seenUUID[p.Type][p.UUID] = struct{}{}
		if p.Location == PackLocationWorld == toWorld {
			seenFolder[p.Type][p.FolderName] = struct{}{}
		}
	}
	for i := range candidates {
		p := &candidates[i]
		if toWorld {
			p.Location = PackLocationWorld
		} else {
			p.Location = PackLocationServer
		}
		if seenUUID[p.Type] == nil {
			seenUUID[p.Type] = map[string]struct{}{}
			seenFolder[p.Type] = map[string]struct{}{}
		}
		if _, ok := seenUUID[p.Type][p.UUID]; ok {
			return fmt.Errorf("a %s pack with UUID %s already exists", p.Type, p.UUID)
		}
		if _, ok := seenFolder[p.Type][p.FolderName]; ok {
			return fmt.Errorf("destination directory already exists: %s", p.FolderName)
		}
		seenUUID[p.Type][p.UUID] = struct{}{}
		seenFolder[p.Type][p.FolderName] = struct{}{}
	}
	return nil
}
