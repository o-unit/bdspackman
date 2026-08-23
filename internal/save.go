package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

// SaveOptions controls save behavior.
type SaveOptions struct {
	ExportPrefix string
	ExportDir    string
}

// SaveWorldPackFiles updates world_behavior_packs.json and
// world_resource_packs.json.
func SaveWorldPackFiles(
	cfg Config,
	packs []Pack,
	opt SaveOptions,
) (string, error) {

	behavior := buildWorldPackList(
		packs,
		PackTypeBehavior,
	)

	resource := buildWorldPackList(
		packs,
		PackTypeResource,
	)

	// Export mode
	if opt.ExportPrefix != "" || opt.ExportDir != "" {
		behaviorDir := filepath.Dir(cfg.WorldBehaviorJSON())
		resourceDir := filepath.Dir(cfg.WorldResourceJSON())

		if opt.ExportDir != "" {
			behaviorDir = opt.ExportDir
			resourceDir = opt.ExportDir

			if err := EnsureDir(opt.ExportDir); err != nil {
				return "", err
			}
		}

		behaviorPath := filepath.Join(
			behaviorDir,
			opt.ExportPrefix+filepath.Base(cfg.WorldBehaviorJSON()),
		)

		resourcePath := filepath.Join(
			resourceDir,
			opt.ExportPrefix+filepath.Base(cfg.WorldResourceJSON()),
		)

		if err := saveWorldPackFilePair(behaviorPath, behavior, resourcePath, resource); err != nil {
			return "", err
		}

		return "Exported.", nil
	}

	if err := saveWorldPackFilePair(cfg.WorldBehaviorJSON(), behavior, cfg.WorldResourceJSON(), resource); err != nil {
		return "", err
	}

	return "Saved.", nil
}

// buildWorldPackList creates world_*_packs.json contents.
func buildWorldPackList(
	packs []Pack,
	packType PackType,
) []WorldPack {

	list := make([]WorldPack, 0)

	for _, pack := range packs {

		if pack.Type != packType {
			continue
		}

		if pack.Status != StatusOn {
			continue
		}

		list = append(list, WorldPack{
			PackID:  pack.UUID,
			Version: pack.Version,
		})
	}

	return list
}

type worldPackFile struct {
	path string
	data []byte
}

// saveWorldPackFilePair replaces both world pack files together. It prepares
// both new files first and restores the original pair if either replacement
// fails.
func saveWorldPackFilePair(behaviorPath string, behavior []WorldPack, resourcePath string, resource []WorldPack) error {
	behaviorData, err := marshalWorldPackJSON(behavior)
	if err != nil {
		return err
	}
	resourceData, err := marshalWorldPackJSON(resource)
	if err != nil {
		return err
	}

	return replaceWorldPackFiles([]worldPackFile{
		{path: behaviorPath, data: behaviorData},
		{path: resourcePath, data: resourceData},
	})
}

func replaceWorldPackFiles(files []worldPackFile) error {
	temps := make([]string, len(files))
	backups := make([]string, len(files))
	hadOriginal := make([]bool, len(files))
	defer func() {
		for _, path := range temps {
			if path != "" {
				_ = os.Remove(path)
			}
		}
		for _, path := range backups {
			if path != "" {
				_ = os.Remove(path)
			}
		}
	}()

	for i, file := range files {
		if err := EnsureDir(filepath.Dir(file.path)); err != nil {
			return err
		}
		tmp, err := os.CreateTemp(filepath.Dir(file.path), "."+filepath.Base(file.path)+".*.tmp")
		if err != nil {
			return fmt.Errorf("create temporary file for %s: %w", file.path, err)
		}
		temps[i] = tmp.Name()
		if _, err := tmp.Write(file.data); err != nil {
			_ = tmp.Close()
			return fmt.Errorf("write temporary file for %s: %w", file.path, err)
		}
		if err := tmp.Chmod(0644); err != nil {
			_ = tmp.Close()
			return fmt.Errorf("chmod temporary file for %s: %w", file.path, err)
		}
		if err := tmp.Close(); err != nil {
			return fmt.Errorf("close temporary file for %s: %w", file.path, err)
		}
	}

	for i, file := range files {
		if _, err := os.Lstat(file.path); err == nil {
			backup, err := os.CreateTemp(filepath.Dir(file.path), "."+filepath.Base(file.path)+".*.bak")
			if err != nil {
				return fmt.Errorf("create backup path for %s: %w", file.path, err)
			}
			backups[i] = backup.Name()
			if err := backup.Close(); err != nil {
				return fmt.Errorf("close backup path for %s: %w", file.path, err)
			}
			if err := os.Rename(file.path, backups[i]); err != nil {
				return rollbackWorldPackFiles(files, backups, hadOriginal, i, err)
			}
			hadOriginal[i] = true
		} else if !os.IsNotExist(err) {
			return rollbackWorldPackFiles(files, backups, hadOriginal, i, err)
		}

		if err := os.Rename(temps[i], file.path); err != nil {
			return rollbackWorldPackFiles(files, backups, hadOriginal, i, err)
		}
		temps[i] = ""
	}

	return nil
}

func rollbackWorldPackFiles(files []worldPackFile, backups []string, hadOriginal []bool, attempted int, cause error) error {
	for i := attempted; i >= 0; i-- {
		if !hadOriginal[i] {
			_ = os.Remove(files[i].path)
			continue
		}
		if err := os.Remove(files[i].path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("replace world pack files: %w (rollback remove %s: %v)", cause, files[i].path, err)
		}
		if err := os.Rename(backups[i], files[i].path); err != nil {
			return fmt.Errorf("replace world pack files: %w (rollback restore %s: %v)", cause, files[i].path, err)
		}
		backups[i] = ""
	}

	return fmt.Errorf("replace world pack files: %w", cause)
}
