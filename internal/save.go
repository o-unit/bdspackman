package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SaveOptions controls save behavior.
type SaveOptions struct {
	DryRun bool
}

// SaveWorldPackFiles updates world_behavior_packs.json and
// world_resource_packs.json.
func SaveWorldPackFiles(
	cfg Config,
	packs []Pack,
	opt SaveOptions,
) error {

	behavior := buildWorldPackList(packs, PackTypeBehavior)
	resource := buildWorldPackList(packs, PackTypeResource)

	if opt.DryRun {

		fmt.Println("=== Dry Run ===")

		fmt.Printf(
			"Behavior Packs : %d\n",
			len(behavior),
		)

		fmt.Printf(
			"Resource Packs : %d\n",
			len(resource),
		)

		return nil
	}

	if err := saveWorldPackJSON(
		cfg.WorldBehaviorJSON(),
		behavior,
	); err != nil {
		return err
	}

	if err := saveWorldPackJSON(
		cfg.WorldResourceJSON(),
		resource,
	); err != nil {
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

	if err := writeJSON(tmp, packs); err != nil {
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

func writeJSON(
	path string,
	v any,
) error {

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf(
			"create %s: %w",
			path,
			err,
		)
	}

	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")

	if err := enc.Encode(v); err != nil {
		return fmt.Errorf(
			"encode %s: %w",
			path,
			err,
		)
	}

	return nil
}
