package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadManifest reads and parses a manifest.json file.
func LoadManifest(path string) (*Manifest, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open manifest: %w", err)
	}
	defer file.Close()

	var manifest Manifest

	if err := json.NewDecoder(file).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("decode manifest: %w", err)
	}

	if err := ValidateManifest(&manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// ValidateManifest validates the required fields of a manifest.
func ValidateManifest(m *Manifest) error {
	if m == nil {
		return fmt.Errorf("manifest is nil")
	}

	if m.Header.Name == "" {
		return fmt.Errorf("manifest.header.name is missing")
	}

	if m.Header.UUID == "" {
		return fmt.Errorf("manifest.header.uuid is missing")
	}

	if len(m.Header.Version) != 3 {
		return fmt.Errorf("manifest.header.version must contain 3 elements")
	}

	return nil
}