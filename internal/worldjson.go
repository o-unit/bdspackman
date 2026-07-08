package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// LoadWorldPacks loads a world_*_packs.json file.
//
// If the file does not exist, an empty slice is returned.
func LoadWorldPacks(path string) ([]WorldPack, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []WorldPack{}, nil
		}
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	var packs []WorldPack

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&packs); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}

	return packs, nil
}

// SaveWorldPacks writes a world_*_packs.json file.
func SaveWorldPacks(path string, packs []WorldPack) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(packs); err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}

	return nil
}

// WorldPackSet returns a lookup table for pack IDs.
func WorldPackSet(packs []WorldPack) map[string]struct{} {
	set := make(map[string]struct{}, len(packs))

	for _, p := range packs {
		set[p.PackID] = struct{}{}
	}

	return set
}
