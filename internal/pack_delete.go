package internal

import (
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

	dst := filepath.Join(
		backupDir,
		filepath.Base(pack.Path)+"-"+date,
	)

	return MovePath(pack.Path, dst, FileOperationOptions{Overwrite: false})
}
