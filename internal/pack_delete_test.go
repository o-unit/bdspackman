package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeletePackRemovesPackDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pack")
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	if err := os.WriteFile(filepath.Join(path, "manifest.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if err := DeletePack(Pack{Path: path}); err != nil {
		t.Fatalf("DeletePack: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("pack directory stat error = %v, want not exist", err)
	}
}

func TestDeletePackRejectsEmptyPath(t *testing.T) {
	if err := DeletePack(Pack{}); err == nil {
		t.Fatal("DeletePack succeeded, want empty-path error")
	}
}
