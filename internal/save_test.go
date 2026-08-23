package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReplaceWorldPackFilesPreparesAllFilesBeforeReplacing(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, "first.json")
	if err := os.WriteFile(first, []byte("original"), 0644); err != nil {
		t.Fatal(err)
	}
	blockedParent := filepath.Join(root, "not-a-directory")
	if err := os.WriteFile(blockedParent, []byte("block"), 0644); err != nil {
		t.Fatal(err)
	}

	err := replaceWorldPackFiles([]worldPackFile{
		{path: first, data: []byte("updated")},
		{path: filepath.Join(blockedParent, "second.json"), data: []byte("updated")},
	})
	if err == nil {
		t.Fatal("replacement succeeded with an invalid second destination")
	}
	data, err := os.ReadFile(first)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "original" {
		t.Fatalf("first file = %q, want original", data)
	}
}

func TestSaveWorldPackFilePairWritesBothFiles(t *testing.T) {
	root := t.TempDir()
	behaviorPath := filepath.Join(root, "world_behavior_packs.json")
	resourcePath := filepath.Join(root, "world_resource_packs.json")

	if err := saveWorldPackFilePair(behaviorPath, []WorldPack{{PackID: "bp", Version: []uint32{1, 0, 0}}}, resourcePath, []WorldPack{{PackID: "rp", Version: []uint32{2, 0, 0}}}); err != nil {
		t.Fatal(err)
	}
	behavior, err := LoadWorldPacks(behaviorPath)
	if err != nil {
		t.Fatal(err)
	}
	resource, err := LoadWorldPacks(resourcePath)
	if err != nil {
		t.Fatal(err)
	}
	if len(behavior) != 1 || behavior[0].PackID != "bp" || len(resource) != 1 || resource[0].PackID != "rp" {
		t.Fatalf("saved files = behavior %#v, resource %#v", behavior, resource)
	}
}
