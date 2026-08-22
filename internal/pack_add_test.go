package internal

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestAddPackReturnsNotImplemented(t *testing.T) {
	inputPath := filepath.Join(t.TempDir(), "pack.mcpack")

	packs, err := AddPack(Config{}, inputPath, false)
	if !errors.Is(err, errAddPackNotImplemented) {
		t.Fatalf("AddPack error = %v, want %v", err, errAddPackNotImplemented)
	}

	if packs != nil {
		t.Fatalf("AddPack packs = %v, want nil", packs)
	}
}
