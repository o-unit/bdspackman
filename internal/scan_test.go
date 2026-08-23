package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsPackEnabledRequiresJSONVersionAtMostPackVersion(t *testing.T) {
	enabled := []WorldPack{{PackID: "pack-uuid", Version: []uint32{1, 2, 3}}}

	if !isPackEnabled(enabled, "pack-uuid", []uint32{1, 2, 3}) {
		t.Fatal("matching version was disabled")
	}
	if !isPackEnabled(enabled, "pack-uuid", []uint32{1, 3, 0}) {
		t.Fatal("newer installed version was disabled")
	}
	if isPackEnabled(enabled, "pack-uuid", []uint32{1, 2, 2}) {
		t.Fatal("older installed version was enabled")
	}
	if isPackEnabled(enabled, "other-uuid", []uint32{9, 9, 9}) {
		t.Fatal("different UUID was enabled")
	}
}

func TestScanPacksEnablesPackWhenInstalledVersionIsNewer(t *testing.T) {
	serverDir := t.TempDir()
	cfg := Config{ServerDir: serverDir, WorldName: "world"}
	packDir := filepath.Join(cfg.ServerBehaviorPackDir(), "pack")
	if err := os.MkdirAll(packDir, 0755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"header":{"name":"pack","uuid":"pack-uuid","version":[1,3,0]}}`
	if err := os.WriteFile(filepath.Join(packDir, "manifest.json"), []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(cfg.WorldDir(), 0755); err != nil {
		t.Fatal(err)
	}
	if err := SaveWorldPacks(cfg.WorldBehaviorJSON(), []WorldPack{{PackID: "pack-uuid", Version: []uint32{1, 2, 3}}}); err != nil {
		t.Fatal(err)
	}

	packs, err := ScanPacks(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(packs) != 1 || packs[0].Status != StatusOn {
		t.Fatalf("scanned packs = %#v, want one enabled pack", packs)
	}

	if err := SaveWorldPacks(cfg.WorldBehaviorJSON(), []WorldPack{{PackID: "pack-uuid", Version: []uint32{1, 4, 0}}}); err != nil {
		t.Fatal(err)
	}
	packs, err = ScanPacks(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(packs) != 1 || packs[0].Status != StatusOff {
		t.Fatalf("scanned packs = %#v, want one disabled pack", packs)
	}
}
