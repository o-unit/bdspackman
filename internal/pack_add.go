package internal

import (
	"errors"
	"fmt"
	"os"
)

var errAddPackNotImplemented = errors.New("add pack is not implemented")

// AddPack imports one or more packs into the selected Server or World location.
func AddPack(cfg Config, inputPath string, toWorld bool) ([]Pack, error) {
	workDir, err := os.MkdirTemp("", "bdspackman-add-*")
	if err != nil {
		return nil, fmt.Errorf("create add-pack work directory: %w", err)
	}
	defer os.RemoveAll(workDir)

	// TODO: Validate inputPath and extract its archive or copy its directory.
	// TODO: Recursively extract nested archives and discover pack manifests.
	// TODO: Detect duplicate UUIDs, choose destination names, and install packs.
	// TODO: Return installed packs using cfg and toWorld to select destinations.
	_ = cfg
	_ = inputPath
	_ = toWorld

	return nil, errAddPackNotImplemented
}
