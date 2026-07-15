package service

import (
	"context"
	"errors"
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

type failingTaskStore struct {
	task *model.UploadTask
}

func (s failingTaskStore) GetByUploadID(ctx context.Context, uploadID string, userID int64) (*model.UploadTask, error) {
	if s.task != nil && s.task.UploadID == uploadID && s.task.UserID == userID {
		return s.task, nil
	}
	return nil, errors.New("task not found")
}

func (s failingTaskStore) UpdateProgressGuarded(ctx context.Context, uploadID string, userID int64, uploadedChunks string, status string, expectedStatus string, expectedVersion int64) error {
	return errors.New("update progress failed")
}

type racingTaskStore struct {
	repo         *db.UploadTaskRepository
	beforeUpdate func(ctx context.Context, uploadID string, userID int64, expectedStatus string, expectedVersion int64) error
}

func (s *racingTaskStore) GetByUploadID(ctx context.Context, uploadID string, userID int64) (*model.UploadTask, error) {
	return s.repo.GetByUploadID(ctx, uploadID, userID)
}

func (s *racingTaskStore) UpdateProgressGuarded(ctx context.Context, uploadID string, userID int64, uploadedChunks string, status string, expectedStatus string, expectedVersion int64) error {
	if s.beforeUpdate != nil {
		beforeUpdate := s.beforeUpdate
		s.beforeUpdate = nil
		if err := beforeUpdate(ctx, uploadID, userID, expectedStatus, expectedVersion); err != nil {
			return err
		}
	}
	return s.repo.UpdateProgressGuarded(ctx, uploadID, userID, uploadedChunks, status, expectedStatus, expectedVersion)
}

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

func TestChunkServiceRollsBackChunkOnProgressError(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	tmpDir := t.TempDir()
	tmpStorage := storage.NewTmpStorage(tmpDir)
	chunkRepo := db.NewUploadChunkRepository(database)
	task := &model.UploadTask{
		UploadID:   "upload-rollback",
		UserID:     7,
		ChunkCount: 2,
		Status:     model.UploadTaskStatusUploading,
	}
	svc := NewChunkService(failingTaskStore{task: task}, chunkRepo, tmpStorage)

	_, err = svc.UploadChunk(context.Background(), ChunkInput{
		UserID:     task.UserID,
		UploadID:   task.UploadID,
		ChunkIndex: 0,
		Body:       strings.NewReader("rollback chunk"),
	})
	if err == nil {
		t.Fatal("expected upload chunk to fail")
	}

	stored, err := chunkRepo.ListByUploadID(context.Background(), task.UploadID)
	if err != nil {
		t.Fatalf("list upload chunks: %v", err)
	}
	if len(stored) != 0 {
		t.Fatalf("expected rollback to remove stored chunks, got %d", len(stored))
	}
	chunkPath := filepath.Join(tmpDir, task.UploadID, "chunk_000000.part")
	if _, err := os.Stat(chunkPath); !os.IsNotExist(err) {
		t.Fatalf("expected chunk file to be removed, got err=%v", err)
	}
}

func TestChunkServiceRollsBackChunkWhenTaskChangesBeforeProgressUpdate(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	tmpDir := t.TempDir()
	taskRepo := db.NewUploadTaskRepository(database)
	chunkRepo := db.NewUploadChunkRepository(database)
	tmpStorage := storage.NewTmpStorage(tmpDir)
	task := &model.UploadTask{
		UploadID:   "upload-race",
		UserID:     9,
		ChunkCount: 2,
		Status:     model.UploadTaskStatusUploading,
		ExpiresAt:  time.Now().Add(time.Hour).UTC(),
	}
	if err := taskRepo.Create(context.Background(), task); err != nil {
		t.Fatalf("create upload task: %v", err)
	}

	tasks := &racingTaskStore{
		repo: taskRepo,
		beforeUpdate: func(ctx context.Context, uploadID string, userID int64, expectedStatus string, expectedVersion int64) error {
			return taskRepo.UpdateProgressGuarded(ctx, uploadID, userID, "", model.UploadTaskStatusCancelled, expectedStatus, expectedVersion)
		},
	}
	svc := NewChunkService(tasks, chunkRepo, tmpStorage)

	_, err = svc.UploadChunk(context.Background(), ChunkInput{
		UserID:     task.UserID,
		UploadID:   task.UploadID,
		ChunkIndex: 0,
		Body:       strings.NewReader("race chunk"),
	})
	if !errors.Is(err, db.ErrUploadTaskStateConflict) {
		t.Fatalf("expected guarded update error, got %v", err)
	}

	stored, err := chunkRepo.ListByUploadID(context.Background(), task.UploadID)
	if err != nil {
		t.Fatalf("list upload chunks: %v", err)
	}
	if len(stored) != 0 {
		t.Fatalf("expected rollback to remove stored chunks, got %d", len(stored))
	}

	chunkPath := filepath.Join(tmpDir, task.UploadID, "chunk_000000.part")
	if _, err := os.Stat(chunkPath); !os.IsNotExist(err) {
		t.Fatalf("expected chunk file to be removed, got err=%v", err)
	}

	reloaded, err := taskRepo.GetByUploadID(context.Background(), task.UploadID, task.UserID)
	if err != nil {
		t.Fatalf("reload upload task: %v", err)
	}
	if reloaded.Status != model.UploadTaskStatusCancelled {
		t.Fatalf("expected task to stay cancelled, got %q", reloaded.Status)
	}
	if reloaded.UploadedChunks != "" {
		t.Fatalf("expected uploaded chunks to stay empty, got %q", reloaded.UploadedChunks)
	}
}
