package internal

import (
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestAddPackInstallsDirectoryIntoWorld(t *testing.T) {
	cfg := addPackConfig(t)
	source := createTestPack(t, t.TempDir(), "example", "data", "bp-uuid")

	packs, err := AddPack(cfg, source, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(packs) != 1 {
		t.Fatalf("installed %d packs, want 1", len(packs))
	}
	if packs[0].Location != PackLocationWorld || packs[0].Type != PackTypeBehavior {
		t.Fatalf("installed pack = %#v", packs[0])
	}
	if _, err := os.Stat(filepath.Join(cfg.WorldBehaviorPackDir(), "example", "manifest.json")); err != nil {
		t.Fatal(err)
	}
}

func TestAddPackInstallsArchiveAndRejectsDuplicateUUID(t *testing.T) {
	cfg := addPackConfig(t)
	root := t.TempDir()
	source := createTestPack(t, root, "resource", "resources", "rp-uuid")
	archive := filepath.Join(t.TempDir(), "resource.mcpack")
	writeArchive(t, archive, source)
	if _, err := AddPack(cfg, archive, false); err != nil {
		t.Fatal(err)
	}
	if _, err := AddPack(cfg, archive, true); err == nil {
		t.Fatal("duplicate UUID was accepted")
	}
}

func TestAddPackRejectsArchivePathTraversal(t *testing.T) {
	cfg := addPackConfig(t)
	archive := filepath.Join(t.TempDir(), "unsafe.zip")
	f, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	entry, err := w.Create("../outside.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write([]byte("unsafe")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := AddPack(cfg, archive, false); err == nil {
		t.Fatal("path traversal archive was accepted")
	}
}

func TestAddPackRejectsMissingPath(t *testing.T) {
	_, err := AddPack(addPackConfig(t), filepath.Join(t.TempDir(), "missing.mcpack"), false)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("error = %v, want not exist", err)
	}
}

func addPackConfig(t *testing.T) Config {
	t.Helper()
	root := t.TempDir()
	cfg := Config{ServerDir: root, WorldName: "world"}
	if err := os.MkdirAll(cfg.WorldDir(), 0755); err != nil {
		t.Fatal(err)
	}
	return cfg
}

func createTestPack(t *testing.T, root, name, moduleType, uuid string) string {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"header":{"name":"` + name + `","uuid":"` + uuid + `","version":[1,0,0]},"modules":[{"type":"` + moduleType + `"}]}`
	if err := os.WriteFile(filepath.Join(path, "manifest.json"), []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeArchive(t *testing.T, archive, source string) {
	t.Helper()
	f, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	err = filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		out, err := w.Create(rel)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = out.Write(data)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}
