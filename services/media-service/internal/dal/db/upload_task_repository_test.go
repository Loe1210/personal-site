package db

import (
	"context"
	"testing"
	"time"

	"github.com/Loe1210/personal-site/services/media-service/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestUploadTaskRepositoryStoresStateAndChunks(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := Migrate(database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	ctx := context.Background()
	taskRepo := NewUploadTaskRepository(database)
	chunkRepo := NewUploadChunkRepository(database)
	expiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)

	task := &model.UploadTask{
		UploadID:   "upload-1",
		UserID:     42,
		BizType:    "article",
		BizID:      "article-9",
		FileName:   "video.mp4",
		FileSize:   10 * 1024 * 1024,
		ChunkSize:  1024 * 1024,
		ChunkCount: 10,
		Status:     model.UploadTaskStatusUploading,
		Sha256:     "file-sha256",
		ExpiresAt:  expiresAt,
	}
	if err := taskRepo.Create(ctx, task); err != nil {
		t.Fatalf("create upload task: %v", err)
	}

	chunk := &model.UploadChunk{
		UploadID:    task.UploadID,
		ChunkIndex:  0,
		Size:        task.ChunkSize,
		Sha256:      "chunk-sha256",
		StoragePath: "tmp/upload-1/0.part",
	}
	if err := chunkRepo.Save(ctx, chunk); err != nil {
		t.Fatalf("save upload chunk: %v", err)
	}
	if err := taskRepo.UpdateProgress(ctx, task.UploadID, task.UserID, "0", model.UploadTaskStatusUploading); err != nil {
		t.Fatalf("update upload task progress: %v", err)
	}

	reloaded, err := taskRepo.GetByUploadID(ctx, task.UploadID, task.UserID)
	if err != nil {
		t.Fatalf("reload upload task: %v", err)
	}
	if reloaded.UploadedChunks != "0" {
		t.Fatalf("expected uploaded chunks to be 0, got %q", reloaded.UploadedChunks)
	}
	if reloaded.Status != model.UploadTaskStatusUploading {
		t.Fatalf("expected status uploading, got %q", reloaded.Status)
	}

	chunks, err := chunkRepo.ListByUploadID(ctx, task.UploadID)
	if err != nil {
		t.Fatalf("list upload chunks: %v", err)
	}
	if len(chunks) != 1 {
		t.Fatalf("expected one stored chunk, got %d", len(chunks))
	}
	if chunks[0].ChunkIndex != 0 || chunks[0].StoragePath != "tmp/upload-1/0.part" {
		t.Fatalf("unexpected stored chunk: %+v", chunks[0])
	}
}
