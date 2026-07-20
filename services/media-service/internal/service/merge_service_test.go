package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/Loe1210/personal-site/services/media-service/internal/model"
)

func TestMergeServiceMergesChunksInIndexOrder(t *testing.T) {
	tmpRoot := t.TempDir()
	finalRoot := t.TempDir()
	uploadID := "merge-1"
	chunks := []model.UploadChunk{
		{UploadID: uploadID, ChunkIndex: 1, StoragePath: uploadID + "/chunk_000001.part", Size: 5},
		{UploadID: uploadID, ChunkIndex: 0, StoragePath: uploadID + "/chunk_000000.part", Size: 6},
	}
	for _, item := range []struct{ path, body string }{{chunks[0].StoragePath, "world"}, {chunks[1].StoragePath, "hello "}} {
		path := filepath.Join(tmpRoot, filepath.FromSlash(item.path))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(item.body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	result, err := NewMergeService(tmpRoot, finalRoot, "/static/uploads/images").Merge(context.Background(), MergeInput{UploadID: uploadID, FileName: "cover.txt", ChunkCount: 2, Chunks: chunks})
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	finalPath := filepath.Join(finalRoot, filepath.FromSlash(result.RelativePath))
	data, err := os.ReadFile(finalPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello world" {
		t.Fatalf("content=%q", data)
	}
	info, err := os.Stat(finalPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got&0o444 != 0o444 {
		t.Fatalf("expected merged file to be readable by nginx, got permissions %o", got)
	}
	sum := sha256.Sum256([]byte("hello world"))
	if result.Sha256 != hex.EncodeToString(sum[:]) {
		t.Fatalf("hash=%q", result.Sha256)
	}
}

func TestMergeServiceRejectsFinalHashMismatch(t *testing.T) {
	tmpRoot := t.TempDir()
	finalRoot := t.TempDir()
	uploadID := "merge-hash"
	storagePath := uploadID + "/chunk_000000.part"
	target := filepath.Join(tmpRoot, filepath.FromSlash(storagePath))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("actual"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := NewMergeService(tmpRoot, finalRoot, "/uploads").Merge(context.Background(), MergeInput{UploadID: uploadID, FileName: "x.bin", ChunkCount: 1, ExpectedSHA256: "deadbeef", Chunks: []model.UploadChunk{{UploadID: uploadID, ChunkIndex: 0, StoragePath: storagePath, Size: 6}}})
	if err == nil {
		t.Fatal("expected final hash mismatch")
	}
}
