package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicWriteFileCreatesParentAndReplacesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "data.txt")

	if err := AtomicWriteFile(path, []byte("first"), 0640); err != nil {
		t.Fatalf("AtomicWriteFile first: %v", err)
	}

	if err := AtomicWriteFile(path, []byte("second"), 0640); err != nil {
		t.Fatalf("AtomicWriteFile second: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if string(data) != "second" {
		t.Fatalf("file content = %q, want %q", string(data), "second")
	}
}

func TestCopyDirCopiesTreeAndRejectsExistingDestination(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")

	if err := os.MkdirAll(filepath.Join(src, "child"), 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	if err := os.WriteFile(filepath.Join(src, "child", "pack.txt"), []byte("pack"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if err := CopyDir(src, dst, FileOperationOptions{}); err != nil {
		t.Fatalf("CopyDir: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dst, "child", "pack.txt"))
	if err != nil {
		t.Fatalf("ReadFile copied file: %v", err)
	}

	if string(data) != "pack" {
		t.Fatalf("copied file content = %q, want %q", string(data), "pack")
	}

	if err := CopyDir(src, dst, FileOperationOptions{}); err == nil {
		t.Fatal("CopyDir existing destination succeeded, want error")
	}
}

func TestMovePathMovesFile(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src.txt")
	dst := filepath.Join(root, "nested", "dst.txt")

	if err := os.WriteFile(src, []byte("move"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if err := MovePath(src, dst, FileOperationOptions{}); err != nil {
		t.Fatalf("MovePath: %v", err)
	}

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatalf("source stat error = %v, want not exist", err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("ReadFile destination: %v", err)
	}

	if string(data) != "move" {
		t.Fatalf("destination content = %q, want %q", string(data), "move")
	}
}

func TestCopyDirRejectsDestinationInsideSource(t *testing.T) {
	src := filepath.Join(t.TempDir(), "src")
	dst := filepath.Join(src, "child")

	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	if err := CopyDir(src, dst, FileOperationOptions{}); err == nil {
		t.Fatal("CopyDir destination inside source succeeded, want error")
	}
}
