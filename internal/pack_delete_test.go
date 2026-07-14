package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeletePackMovesPackToDelpacks(t *testing.T) {
	serverDir := t.TempDir()

	path := filepath.Join(serverDir, "behavior_packs", "pack")
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	if err := os.WriteFile(filepath.Join(path, "manifest.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg := Config{
		ServerDir: serverDir,
	}

	pack := Pack{
		Path:     path,
		Location: PackLocationServer,
		Type:     PackTypeBehavior,
	}

	if err := DeletePack(cfg, pack); err != nil {
		t.Fatalf("DeletePack: %v", err)
	}

	// 元ディレクトリは消えていること
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("pack directory stat error = %v, want not exist", err)
	}

	// delpacksへ移動されていること
	backupRoot := filepath.Join(
		serverDir,
		"delpacks",
		"server",
		"behavior_packs",
	)

	entries, err := os.ReadDir(backupRoot)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("backup count = %d, want 1", len(entries))
	}
}

func TestDeletePackRejectsEmptyPath(t *testing.T) {
	cfg := Config{
		ServerDir: t.TempDir(),
	}

	if err := DeletePack(cfg, Pack{}); err == nil {
		t.Fatal("DeletePack succeeded, want empty-path error")
	}
}

func TestDeletePackForceDelete(t *testing.T) {
	serverDir := t.TempDir()

	path := filepath.Join(serverDir, "behavior_packs", "pack")
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	cfg := Config{
		ServerDir:   serverDir,
		ForceDelete: true,
	}

	pack := Pack{
		Path: path,
	}

	if err := DeletePack(cfg, pack); err != nil {
		t.Fatalf("DeletePack: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("pack directory stat error = %v, want not exist", err)
	}

	backupRoot := filepath.Join(serverDir, "delpacks")
	if _, err := os.Stat(backupRoot); !os.IsNotExist(err) {
		t.Fatal("delpacks directory should not be created")
	}
}
