package internal

import (
	"fmt"
	"path/filepath"
	"strings"
)

// RenamePackDirectory renames a pack directory without changing its location.
func RenamePackDirectory(pack *Pack, newName string) error {
	if pack == nil {
		return fmt.Errorf("pack is nil")
	}

	name, err := validatePackDirectoryName(newName)
	if err != nil {
		return err
	}

	if pack.Path == "" {
		return fmt.Errorf("pack path is empty")
	}

	if pack.FolderName == name {
		return nil
	}

	targetPath := filepath.Join(filepath.Dir(pack.Path), name)
	if err := MovePath(
		pack.Path,
		targetPath,
		FileOperationOptions{},
	); err != nil {
		return err
	}

	pack.FolderName = name
	pack.Path = targetPath
	pack.ManifestPath = filepath.Join(targetPath, "manifest.json")

	return nil
}

func validatePackDirectoryName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("directory name is empty")
	}

	if name == "." || name == ".." {
		return "", fmt.Errorf("invalid directory name: %s", name)
	}

	if filepath.Base(name) != name {
		return "", fmt.Errorf("directory name must not contain path separators: %s", name)
	}

	return name, nil
}
