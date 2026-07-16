package service

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
func TestChunkServiceRetryKeepsPreviousChunkWhenProgressUpdateConflicts(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	ctx := context.Background()
	tmpDir := t.TempDir()
	taskRepo := db.NewUploadTaskRepository(database)
	chunkRepo := db.NewUploadChunkRepository(database)
	tmpStorage := storage.NewTmpStorage(tmpDir)
	task := &model.UploadTask{
		UploadID:   "upload-retry-conflict",
		UserID:     11,
		ChunkCount: 2,
		Status:     model.UploadTaskStatusUploading,
		ExpiresAt:  time.Now().Add(time.Hour).UTC(),
	}
	if err := taskRepo.Create(ctx, task); err != nil {
		t.Fatalf("create upload task: %v", err)
	}

	svc := NewChunkService(taskRepo, chunkRepo, tmpStorage)
	firstChunk, err := svc.UploadChunk(ctx, ChunkInput{
		UserID:     task.UserID,
		UploadID:   task.UploadID,
		ChunkIndex: 0,
		Body:       strings.NewReader("original chunk"),
	})
	if err != nil {
		t.Fatalf("upload original chunk: %v", err)
	}

	tasks := &racingTaskStore{
		repo: taskRepo,
		beforeUpdate: func(ctx context.Context, uploadID string, userID int64, expectedStatus string, expectedVersion int64) error {
			return taskRepo.UpdateProgressGuarded(ctx, uploadID, userID, "0", model.UploadTaskStatusCancelled, expectedStatus, expectedVersion)
		},
	}
	svc = NewChunkService(tasks, chunkRepo, tmpStorage)
	_, err = svc.UploadChunk(ctx, ChunkInput{
		UserID:     task.UserID,
		UploadID:   task.UploadID,
		ChunkIndex: 0,
		Body:       strings.NewReader("replacement chunk"),
	})
	if !errors.Is(err, db.ErrUploadTaskStateConflict) {
		t.Fatalf("expected guarded update conflict, got %v", err)
	}

	stored, err := chunkRepo.ListByUploadID(ctx, task.UploadID)
	if err != nil {
		t.Fatalf("list upload chunks: %v", err)
	}
	if len(stored) != 1 {
		t.Fatalf("expected previous chunk row to remain, got %d rows", len(stored))
	}
	if stored[0].ChunkIndex != firstChunk.ChunkIndex || stored[0].StoragePath != firstChunk.StoragePath || stored[0].Sha256 != firstChunk.Sha256 {
		t.Fatalf("expected previous chunk row to be restored, got %+v want %+v", stored[0], firstChunk)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, filepath.FromSlash(firstChunk.StoragePath)))
	if err != nil {
		t.Fatalf("read restored chunk file: %v", err)
	}
	if string(data) != "original chunk" {
		t.Fatalf("expected original chunk content to remain, got %q", string(data))
	}

	reloaded, err := taskRepo.GetByUploadID(ctx, task.UploadID, task.UserID)
	if err != nil {
		t.Fatalf("reload upload task: %v", err)
	}
	if reloaded.Status != model.UploadTaskStatusCancelled {
		t.Fatalf("expected task to stay cancelled, got %q", reloaded.Status)
	}
	if reloaded.UploadedChunks != "0" {
		t.Fatalf("expected uploaded chunks to keep previous progress, got %q", reloaded.UploadedChunks)
	}
}

type blockingProgressStore struct {
	repo    *db.UploadTaskRepository
	entered chan struct{}
	release chan struct{}
	once    sync.Once
}

func (s *blockingProgressStore) GetByUploadID(ctx context.Context, uploadID string, userID int64) (*model.UploadTask, error) {
	return s.repo.GetByUploadID(ctx, uploadID, userID)
}

func (s *blockingProgressStore) UpdateProgressGuarded(ctx context.Context, uploadID string, userID int64, uploadedChunks string, status string, expectedStatus string, expectedVersion int64) error {
	first := false
	s.once.Do(func() { first = true })
	if first {
		close(s.entered)
		<-s.release
		if err := s.repo.UpdateProgressGuarded(ctx, uploadID, userID, uploadedChunks, status, expectedStatus, expectedVersion); err != nil {
			return err
		}
	}
	return s.repo.UpdateProgressGuarded(ctx, uploadID, userID, uploadedChunks, status, expectedStatus, expectedVersion)
}

type observingChunkStorage struct {
	store      *storage.TmpStorage
	mu         sync.Mutex
	saveCount  int
	secondSave chan struct{}
	once       sync.Once
}

func (s *observingChunkStorage) SaveChunk(uploadID string, chunkIndex int, content io.Reader) (string, int64, string, error) {
	s.mu.Lock()
	s.saveCount++
	isSecond := s.saveCount == 2
	s.mu.Unlock()
	if isSecond {
		s.once.Do(func() { close(s.secondSave) })
	}
	return s.store.SaveChunk(uploadID, chunkIndex, content)
}

func (s *observingChunkStorage) BackupChunk(storagePath string) (string, bool, error) {
	return s.store.BackupChunk(storagePath)
}

func (s *observingChunkStorage) RestoreChunk(storagePath string, backupPath string) error {
	return s.store.RestoreChunk(storagePath, backupPath)
}

func (s *observingChunkStorage) DiscardChunkBackup(backupPath string) error {
	return s.store.DiscardChunkBackup(backupPath)
}

func (s *observingChunkStorage) RemoveChunk(storagePath string) error {
	return s.store.RemoveChunk(storagePath)
}

func TestChunkServiceSerializesConcurrentRetriesForSameChunk(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	ctx := context.Background()
	tmpDir := t.TempDir()
	taskRepo := db.NewUploadTaskRepository(database)
	chunkRepo := db.NewUploadChunkRepository(database)
	task := &model.UploadTask{
		UploadID:   "upload-concurrent-retry",
		UserID:     12,
		ChunkCount: 2,
		Status:     model.UploadTaskStatusUploading,
		ExpiresAt:  time.Now().Add(time.Hour).UTC(),
	}
	if err := taskRepo.Create(ctx, task); err != nil {
		t.Fatalf("create upload task: %v", err)
	}

	tmpStorage := storage.NewTmpStorage(tmpDir)
	seedService := NewChunkService(taskRepo, chunkRepo, tmpStorage)
	if _, err := seedService.UploadChunk(ctx, ChunkInput{
		UserID: task.UserID, UploadID: task.UploadID, ChunkIndex: 0, Body: strings.NewReader("original chunk"),
	}); err != nil {
		t.Fatalf("seed original chunk: %v", err)
	}

	progressStore := &blockingProgressStore{repo: taskRepo, entered: make(chan struct{}), release: make(chan struct{})}
	observedStorage := &observingChunkStorage{store: tmpStorage, secondSave: make(chan struct{})}
	svc := NewChunkService(progressStore, chunkRepo, observedStorage)

	firstResult := make(chan error, 1)
	go func() {
		_, err := svc.UploadChunk(ctx, ChunkInput{
			UserID: task.UserID, UploadID: task.UploadID, ChunkIndex: 0, Body: strings.NewReader("first retry"),
		})
		firstResult <- err
	}()
	<-progressStore.entered

	secondResult := make(chan error, 1)
	go func() {
		_, err := svc.UploadChunk(ctx, ChunkInput{
			UserID: task.UserID, UploadID: task.UploadID, ChunkIndex: 0, Body: strings.NewReader("second retry"),
		})
		secondResult <- err
	}()

	select {
	case <-observedStorage.secondSave:
		t.Fatal("second retry wrote the chunk before the first retry completed")
	case <-time.After(100 * time.Millisecond):
	}

	close(progressStore.release)
	if err := <-firstResult; !errors.Is(err, db.ErrUploadTaskStateConflict) {
		t.Fatalf("expected first retry to lose the version race, got %v", err)
	}
	if err := <-secondResult; err != nil {
		t.Fatalf("second retry should complete: %v", err)
	}

	stored, err := chunkRepo.ListByUploadID(ctx, task.UploadID)
	if err != nil {
		t.Fatalf("list chunks: %v", err)
	}
	if len(stored) != 1 {
		t.Fatalf("expected one chunk, got %d", len(stored))
	}
	data, err := os.ReadFile(filepath.Join(tmpDir, filepath.FromSlash(stored[0].StoragePath)))
	if err != nil {
		t.Fatalf("read final chunk: %v", err)
	}
	if string(data) != "second retry" {
		t.Fatalf("expected second retry to own final chunk, got %q", string(data))
	}
}
