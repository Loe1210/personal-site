package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTmpStorageWritesChunkToTmpPath(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewTmpStorage(tmpDir)

	storagePath, size, digest, err := store.SaveChunk("upload-1", 2, strings.NewReader("hello chunk"))
	if err != nil {
		t.Fatalf("save chunk: %v", err)
	}
	if storagePath != "upload-1/chunk_000002.part" {
		t.Fatalf("unexpected storage path: %q", storagePath)
	}
	if size != int64(len("hello chunk")) {
		t.Fatalf("unexpected size: %d", size)
	}
	if digest == "" {
		t.Fatal("expected digest to be populated")
	}
	data, err := os.ReadFile(filepath.Join(tmpDir, filepath.FromSlash(storagePath)))
	if err != nil {
		t.Fatalf("read chunk: %v", err)
	}
	if string(data) != "hello chunk" {
		t.Fatalf("unexpected chunk content: %q", string(data))
	}
}

func TestTmpStorageReplacesExistingChunk(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewTmpStorage(tmpDir)

	storagePath, _, _, err := store.SaveChunk("upload-1", 2, strings.NewReader("first"))
	if err != nil {
		t.Fatalf("save first chunk: %v", err)
	}
	secondPath, _, _, err := store.SaveChunk("upload-1", 2, strings.NewReader("second"))
	if err != nil {
		t.Fatalf("save second chunk: %v", err)
	}
	if secondPath != storagePath {
		t.Fatalf("expected retry to reuse same path, got %q and %q", storagePath, secondPath)
	}
	data, err := os.ReadFile(filepath.Join(tmpDir, filepath.FromSlash(storagePath)))
	if err != nil {
		t.Fatalf("read chunk: %v", err)
	}
	if string(data) != "second" {
		t.Fatalf("expected retry to replace chunk content, got %q", string(data))
	}
}

func TestTmpStorageRejectsUploadIDPathTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewTmpStorage(tmpDir)

	if _, _, _, err := store.SaveChunk("../escape", 0, strings.NewReader("bad")); err == nil {
		t.Fatal("expected traversal upload id to be rejected")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "..", "escape")); !os.IsNotExist(err) {
		t.Fatalf("expected no escaped directory, got %v", err)
	}
}
