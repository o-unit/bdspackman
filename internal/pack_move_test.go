package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMovePackMovesServerBehaviorPackToWorld(t *testing.T) {
	serverDir := t.TempDir()
	cfg := Config{
		ServerDir: serverDir,
		WorldName: "world",
	}

	src := filepath.Join(cfg.ServerBehaviorPackDir(), "pack")
	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatalf("MkdirAll source: %v", err)
	}

	if err := os.WriteFile(filepath.Join(src, "manifest.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("WriteFile manifest: %v", err)
	}

	pack := Pack{
		Type:       PackTypeBehavior,
		Location:   PackLocationServer,
		FolderName: "pack",
		Path:       src,
	}

	if err := MovePack(cfg, &pack); err != nil {
		t.Fatalf("MovePack: %v", err)
	}

	wantPath := filepath.Join(cfg.WorldBehaviorPackDir(), "pack")
	if pack.Location != PackLocationWorld {
		t.Fatalf("pack location = %v, want %v", pack.Location, PackLocationWorld)
	}

	if pack.Path != wantPath {
		t.Fatalf("pack path = %q, want %q", pack.Path, wantPath)
	}

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatalf("source stat error = %v, want not exist", err)
	}

	if _, err := os.Stat(filepath.Join(wantPath, "manifest.json")); err != nil {
		t.Fatalf("destination manifest stat: %v", err)
	}
}

func TestMovePackMovesWorldResourcePackToServer(t *testing.T) {
	serverDir := t.TempDir()
	cfg := Config{
		ServerDir: serverDir,
		WorldName: "world",
	}

	src := filepath.Join(cfg.WorldResourcePackDir(), "pack")
	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatalf("MkdirAll source: %v", err)
	}

	pack := Pack{
		Type:       PackTypeResource,
		Location:   PackLocationWorld,
		FolderName: "pack",
		Path:       src,
	}

	if err := MovePack(cfg, &pack); err != nil {
		t.Fatalf("MovePack: %v", err)
	}

	wantPath := filepath.Join(cfg.ServerResourcePackDir(), "pack")
	if pack.Location != PackLocationServer {
		t.Fatalf("pack location = %v, want %v", pack.Location, PackLocationServer)
	}

	if pack.Path != wantPath {
		t.Fatalf("pack path = %q, want %q", pack.Path, wantPath)
	}
}

func TestMovePackRejectsExistingDestination(t *testing.T) {
	serverDir := t.TempDir()
	cfg := Config{
		ServerDir: serverDir,
		WorldName: "world",
	}

	src := filepath.Join(cfg.ServerBehaviorPackDir(), "pack")
	dst := filepath.Join(cfg.WorldBehaviorPackDir(), "pack")

	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatalf("MkdirAll source: %v", err)
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		t.Fatalf("MkdirAll destination: %v", err)
	}

	pack := Pack{
		Type:       PackTypeBehavior,
		Location:   PackLocationServer,
		FolderName: "pack",
		Path:       src,
	}

	if err := MovePack(cfg, &pack); err == nil {
		t.Fatal("MovePack succeeded, want destination-exists error")
	}
}
