package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DeletePack deletes a pack.
func DeletePack(cfg Config, pack Pack) error {
	// --force-delete
	if cfg.ForceDelete {
		return RemovePath(pack.Path)
	}

	root := cfg.ServerDir

	location := "server"
	if pack.Location == PackLocationWorld {
		location = "world"
	}

	kind := "behavior_packs"
	if pack.Type == PackTypeResource {
		kind = "resource_packs"
	}

	backupDir := filepath.Join(
		root,
		"delpacks",
		location,
		kind,
	)

	if err := EnsureDir(backupDir); err != nil {
		return err
	}

	date := time.Now().Format("060102")

	dst, err := nextDeleteBackupPath(backupDir, filepath.Base(pack.Path), date)
	if err != nil {
		return err
	}

	return MovePath(pack.Path, dst, FileOperationOptions{Overwrite: false})
}

func nextDeleteBackupPath(backupDir, name, date string) (string, error) {
	for sequence := 0; ; sequence++ {
		path := filepath.Join(backupDir, fmt.Sprintf("%s-%s-%02d", name, date, sequence))
		_, err := os.Lstat(path)
		if os.IsNotExist(err) {
			return path, nil
		}
		if err != nil {
			return "", err
		}
	}
}
