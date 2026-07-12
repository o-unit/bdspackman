package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// LoadWorldPacks loads a world_*_packs.json file.
//
// If the file does not exist, an empty slice is returned.
func LoadWorldPacks(path string) ([]WorldPack, error) {

	data, err := os.ReadFile(path)
	if err != nil {

		if errors.Is(err, os.ErrNotExist) {
			return []WorldPack{}, nil
		}

		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	return unmarshalWorldPackJSON(data, path)
}

// SaveWorldPacks writes a world_*_packs.json file.
func SaveWorldPacks(
	path string,
	packs []WorldPack,
) error {

	data, err := marshalWorldPackJSON(packs)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	return nil
}

// marshalWorldPackJSON converts WorldPack to formatted JSON.
func marshalWorldPackJSON(
	packs []WorldPack,
) ([]byte, error) {

	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "    ")

	if err := enc.Encode(packs); err != nil {
		return nil, err
	}

	data := buf.Bytes()

	// version を1行化
	data = bytes.ReplaceAll(
		data,
		[]byte("[\n            "),
		[]byte("[ "),
	)

	data = bytes.ReplaceAll(
		data,
		[]byte(",\n            "),
		[]byte(", "),
	)

	data = bytes.ReplaceAll(
		data,
		[]byte("\n        ]"),
		[]byte(" ]"),
	)

	return data, nil
}

// unmarshalWorldPackJSON converts JSON to WorldPack.
func unmarshalWorldPackJSON(
	data []byte,
	path string,
) ([]WorldPack, error) {

	var packs []WorldPack

	if err := json.Unmarshal(data, &packs); err != nil {
		return nil, fmt.Errorf(
			"decode %s: %w",
			path,
			err,
		)
	}

	return packs, nil
}

// WorldPackSet returns a lookup table for pack IDs.
func WorldPackSet(
	packs []WorldPack,
) map[string]struct{} {

	set := make(
		map[string]struct{},
		len(packs),
	)

	for _, p := range packs {
		set[p.PackID] = struct{}{}
	}

	return set
}
