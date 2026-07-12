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

			if err := os.MkdirAll(opt.ExportDir, 0755); err != nil {
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

		if err := exportWorldPackJSON(
			behaviorPath,
			behavior,
		); err != nil {
			return "", err
		}

		if err := exportWorldPackJSON(
			resourcePath,
			resource,
		); err != nil {
			return "", err
		}

		return "Exported.", nil
	}

	// Normal save
	if err := saveWorldPackJSON(
		cfg.WorldBehaviorJSON(),
		behavior,
	); err != nil {
		return "", err
	}

	if err := saveWorldPackJSON(
		cfg.WorldResourceJSON(),
		resource,
	); err != nil {
		return "", err
	}

	return "Saved.", nil
}

// exportWorldPackJSON writes one world_*_packs.json directly.
func exportWorldPackJSON(
	path string,
	packs []WorldPack,
) error {

	if err := SaveWorldPacks(path, packs); err != nil {
		return err
	}

	return nil
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

// saveWorldPackJSON writes one world_*_packs.json.
func saveWorldPackJSON(
	path string,
	packs []WorldPack,
) error {

	tmp := path + ".tmp"

	if err := SaveWorldPacks(tmp, packs); err != nil {
		return err
	}

	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf(
			"rename %s: %w",
			filepath.Base(path),
			err,
		)
	}

	return nil
}

// writeJSON writes already-generated JSON data.
func writeJSON(
	path string,
	data []byte,
) error {

	if err := os.WriteFile(
		path,
		data,
		0644,
	); err != nil {

		return fmt.Errorf(
			"write %s: %w",
			path,
			err,
		)
	}

	return nil
}
