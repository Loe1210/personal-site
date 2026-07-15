package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	db "github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
	"github.com/Loe1210/personal-site/services/media-service/internal/dal/storage"
	"github.com/Loe1210/personal-site/services/media-service/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestChunkServiceWritesChunkToTmpPath(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	tmpDir := t.TempDir()
	tmpStorage := storage.NewTmpStorage(tmpDir)
	taskRepo := db.NewUploadTaskRepository(database)
	chunkRepo := db.NewUploadChunkRepository(database)
	svc := NewChunkService(taskRepo, chunkRepo, tmpStorage)

	ctx := context.Background()
	task := &model.UploadTask{
		UploadID:   "upload-1",
		UserID:     42,
		BizType:    "article",
		BizID:      "article-9",
		FileName:   "video.mp4",
		FileSize:   8 * 1024 * 1024,
		ChunkSize:  4 * 1024 * 1024,
		ChunkCount: 2,
		Status:     model.UploadTaskStatusUploading,
		ExpiresAt:  time.Now().Add(time.Hour).UTC(),
	}
	if err := taskRepo.Create(ctx, task); err != nil {
		t.Fatalf("create upload task: %v", err)
	}

	chunk, err := svc.UploadChunk(ctx, ChunkInput{
		UserID:     task.UserID,
		UploadID:   task.UploadID,
		ChunkIndex: 1,
		Body:       strings.NewReader("hello chunk"),
	})
	if err != nil {
		t.Fatalf("upload chunk: %v", err)
	}
	if chunk.StoragePath != "upload-1/chunk_000001.part" {
		t.Fatalf("unexpected storage path: %q", chunk.StoragePath)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, filepath.FromSlash(chunk.StoragePath)))
	if err != nil {
		t.Fatalf("read chunk file: %v", err)
	}
	if string(data) != "hello chunk" {
		t.Fatalf("unexpected chunk content: %q", string(data))
	}

	reloaded, err := taskRepo.GetByUploadID(ctx, task.UploadID, task.UserID)
	if err != nil {
		t.Fatalf("reload upload task: %v", err)
	}
	if reloaded.UploadedChunks != "1" {
		t.Fatalf("expected uploaded chunks to be 1, got %q", reloaded.UploadedChunks)
	}
	if reloaded.Status != model.UploadTaskStatusUploading {
		t.Fatalf("expected status uploading, got %q", reloaded.Status)
	}

	stored, err := chunkRepo.ListByUploadID(ctx, task.UploadID)
	if err != nil {
		t.Fatalf("list upload chunks: %v", err)
	}
	if len(stored) != 1 {
		t.Fatalf("expected one stored chunk, got %d", len(stored))
	}
	if stored[0].ChunkIndex != 1 || stored[0].StoragePath != chunk.StoragePath {
		t.Fatalf("unexpected stored chunk: %+v", stored[0])
	}
}
