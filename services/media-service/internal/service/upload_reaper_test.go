package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
	"github.com/Loe1210/personal-site/services/media-service/internal/dal/storage"
	"github.com/Loe1210/personal-site/services/media-service/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestUploadReaperDeletesExpiredTmpFiles(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	ctx := context.Background()
	tmpRoot := t.TempDir()
	taskRepo := db.NewUploadTaskRepository(database)
	chunkRepo := db.NewUploadChunkRepository(database)
	store := storage.NewTmpStorage(tmpRoot)
	task := &model.UploadTask{
		UploadID:   "expired-upload",
		UserID:     8,
		FileName:   "old.bin",
		FileSize:   5,
		ChunkSize:  5,
		ChunkCount: 1,
		Status:     model.UploadTaskStatusUploading,
		ExpiresAt:  time.Now().Add(-time.Hour).UTC(),
	}
	if err := taskRepo.Create(ctx, task); err != nil {
		t.Fatalf("create upload task: %v", err)
	}
	if err := chunkRepo.Save(ctx, &model.UploadChunk{UploadID: task.UploadID, ChunkIndex: 0, Size: 5, StoragePath: task.UploadID + "/chunk_000000.part"}); err != nil {
		t.Fatalf("save chunk: %v", err)
	}
	chunkPath := filepath.Join(tmpRoot, task.UploadID, "chunk_000000.part")
	if err := os.MkdirAll(filepath.Dir(chunkPath), 0o755); err != nil {
		t.Fatalf("create tmp dir: %v", err)
	}
	if err := os.WriteFile(chunkPath, []byte("stale"), 0o644); err != nil {
		t.Fatalf("write tmp chunk: %v", err)
	}

	reaper := NewUploadReaper(taskRepo, store, 100)
	deleted, err := reaper.RunOnce(ctx, time.Now().UTC())
	if err != nil {
		t.Fatalf("run reaper: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("expected one deleted task, got %d", deleted)
	}
	if _, err := os.Stat(filepath.Join(tmpRoot, task.UploadID)); !os.IsNotExist(err) {
		t.Fatalf("expected tmp upload directory to be removed, got %v", err)
	}
	if _, err := taskRepo.GetByUploadID(ctx, task.UploadID, task.UserID); err == nil {
		t.Fatal("expected expired task to be deleted")
	}
}
