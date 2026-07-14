package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenamePackDirectoryRenamesDirectoryAndUpdatesPack(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "old_name")
	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	if err := os.WriteFile(filepath.Join(src, "manifest.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	pack := Pack{
		FolderName: "old_name",
		Path:       src,
	}

	if err := RenamePackDirectory(&pack, "new_name"); err != nil {
		t.Fatalf("RenamePackDirectory: %v", err)
	}

	wantPath := filepath.Join(root, "new_name")
	if pack.FolderName != "new_name" {
		t.Fatalf("folder name = %q, want %q", pack.FolderName, "new_name")
	}

	if pack.Path != wantPath {
		t.Fatalf("path = %q, want %q", pack.Path, wantPath)
	}

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatalf("old directory stat error = %v, want not exist", err)
	}

	if _, err := os.Stat(filepath.Join(wantPath, "manifest.json")); err != nil {
		t.Fatalf("new manifest stat: %v", err)
	}
}

func TestRenamePackDirectoryRejectsPathSeparator(t *testing.T) {
	pack := Pack{
		FolderName: "pack",
		Path:       filepath.Join(t.TempDir(), "pack"),
	}

	if err := RenamePackDirectory(&pack, filepath.Join("nested", "pack")); err == nil {
		t.Fatal("RenamePackDirectory succeeded, want separator error")
	}
}

func TestRenamePackDirectoryRejectsExistingDestination(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "pack")
	dst := filepath.Join(root, "existing")

	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatalf("MkdirAll source: %v", err)
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		t.Fatalf("MkdirAll destination: %v", err)
	}

	pack := Pack{
		FolderName: "pack",
		Path:       src,
	}

	if err := RenamePackDirectory(&pack, "existing"); err == nil {
		t.Fatal("RenamePackDirectory succeeded, want destination-exists error")
	}
}
