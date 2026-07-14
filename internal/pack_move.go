package internal

import (
	"fmt"
	"path/filepath"
)

// PackMoveTarget returns the location and directory a pack moves to.
func PackMoveTarget(cfg Config, pack Pack) (PackLocation, string) {
	location := PackLocationWorld
	if pack.Location == PackLocationWorld {
		location = PackLocationServer
	}

	return location, packDirectoryForLocation(cfg, pack.Type, location)
}

// MovePack moves a pack directory between the server and world pack locations.
func MovePack(cfg Config, pack *Pack) error {
	if pack == nil {
		return fmt.Errorf("pack is nil")
	}

	targetLocation, targetDir := PackMoveTarget(cfg, *pack)
	targetPath := filepath.Join(targetDir, pack.FolderName)

	if err := MovePath(
		pack.Path,
		targetPath,
		FileOperationOptions{},
	); err != nil {
		return err
	}

	pack.Location = targetLocation
	pack.Path = targetPath
	pack.ManifestPath = filepath.Join(targetPath, "manifest.json")

	return nil
}

func packDirectoryForLocation(
	cfg Config,
	packType PackType,
	location PackLocation,
) string {
	switch {
	case packType == PackTypeBehavior && location == PackLocationServer:
		return cfg.ServerBehaviorPackDir()
	case packType == PackTypeResource && location == PackLocationServer:
		return cfg.ServerResourcePackDir()
	case packType == PackTypeBehavior && location == PackLocationWorld:
		return cfg.WorldBehaviorPackDir()
	default:
		return cfg.WorldResourcePackDir()
	}
}
