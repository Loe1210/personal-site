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
