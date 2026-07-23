package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
	"github.com/Loe1210/personal-site/services/media-service/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type capturingUploadTaskStore struct {
	task            *model.UploadTask
	updateUploadID  string
	updateUserID    int64
	updateChunks    string
	updateStatus    string
	expectedStatus  string
	expectedVersion int64
	updateErr       error
}

func (s *capturingUploadTaskStore) Create(ctx context.Context, task *model.UploadTask) error {
	return nil
}

func (s *capturingUploadTaskStore) GetByUploadID(ctx context.Context, uploadID string, userID int64) (*model.UploadTask, error) {
	if s.task == nil {
		return nil, errors.New("task not found")
	}
	return s.task, nil
}

func (s *capturingUploadTaskStore) UpdateProgressGuarded(ctx context.Context, uploadID string, userID int64, uploadedChunks string, status string, expectedStatus string, expectedVersion int64) error {
	s.updateUploadID = uploadID
	s.updateUserID = userID
	s.updateChunks = uploadedChunks
	s.updateStatus = status
	s.expectedStatus = expectedStatus
	s.expectedVersion = expectedVersion
	return s.updateErr
}

func TestInitUploadRejectsTooLargeFile(t *testing.T) {
	svc := NewUploadTaskService(&configs.UploadConfig{MaxImageSizeMB: 1}, db.NewUploadTaskRepository(nil), db.NewUploadChunkRepository(nil))

	_, err := svc.InitUpload(context.Background(), InitInput{
		UserID:   1,
		FileName: "large.png",
		FileSize: 2 * 1024 * 1024,
	})
	if err == nil {
		t.Fatal("expected too large file to be rejected")
	}
	if !strings.Contains(err.Error(), "too large") {
		t.Fatalf("expected too large error, got %v", err)
	}
}

func TestInitUploadRejectsUnsupportedContentType(t *testing.T) {
	svc := NewUploadTaskService(&configs.UploadConfig{}, db.NewUploadTaskRepository(nil), db.NewUploadChunkRepository(nil))

	_, err := svc.InitUpload(context.Background(), InitInput{
		UserID:      1,
		FileName:    "note.txt",
		FileSize:    10,
		ContentType: "text/plain",
	})
	if err == nil {
		t.Fatal("expected unsupported content type to be rejected")
	}
	if !strings.Contains(err.Error(), "image uploads") {
		t.Fatalf("expected image-only error, got %v", err)
	}
}
func TestCancelUploadUsesTaskStatusAndVersionGuard(t *testing.T) {
	store := &capturingUploadTaskStore{
		task: &model.UploadTask{
			UploadID:       "upload-1",
			UserID:         7,
			UploadedChunks: "0,1",
			Status:         model.UploadTaskStatusUploading,
			Version:        12,
		},
	}
	svc := NewUploadTaskService(&configs.UploadConfig{}, store, db.NewUploadChunkRepository(nil))

	if err := svc.CancelUpload(context.Background(), "upload-1", 7); err != nil {
		t.Fatalf("cancel upload: %v", err)
	}
	if store.updateUploadID != "upload-1" || store.updateUserID != 7 {
		t.Fatalf("unexpected update target: upload_id=%q user_id=%d", store.updateUploadID, store.updateUserID)
	}
	if store.updateChunks != "0,1" {
		t.Fatalf("expected uploaded chunks to be preserved, got %q", store.updateChunks)
	}
	if store.updateStatus != model.UploadTaskStatusCancelled {
		t.Fatalf("expected cancelled status, got %q", store.updateStatus)
	}
	if store.expectedStatus != model.UploadTaskStatusUploading {
		t.Fatalf("expected guard status uploading, got %q", store.expectedStatus)
	}
	if store.expectedVersion != 12 {
		t.Fatalf("expected guard version 12, got %d", store.expectedVersion)
	}
}

func TestCompleteUploadRequiresConfiguredMergePipeline(t *testing.T) {
	store := &capturingUploadTaskStore{task: &model.UploadTask{UploadID: "complete-1", UserID: 7, Status: model.UploadTaskStatusUploading, Version: 1}}
	svc := NewUploadTaskService(&configs.UploadConfig{}, store, db.NewUploadChunkRepository(nil))
	if _, err := svc.CompleteUpload(context.Background(), "complete-1", 7); err == nil {
		t.Fatal("expected completion without merge pipeline to fail")
	}
}

func TestCompleteUploadMergesChunksAndCreatesImageRecord(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	ctx := context.Background()
	tmpRoot := t.TempDir()
	finalRoot := t.TempDir()
	taskRepo := db.NewUploadTaskRepository(database)
	chunkRepo := db.NewUploadChunkRepository(database)
	fileRepo := db.NewFileRepository(database)

	imageBytes := testPNGBytes(t)
	sum := sha256.Sum256(imageBytes)
	task := &model.UploadTask{
		UploadID:       "complete-image",
		UserID:         9,
		BizType:        "article",
		BizID:          "cover-1",
		FileName:       "cover.png",
		FileSize:       int64(len(imageBytes)),
		ContentType:    "image/png",
		ChunkSize:      int64(len(imageBytes)),
		ChunkCount:     1,
		UploadedChunks: "0",
		Status:         model.UploadTaskStatusUploading,
		Sha256:         hex.EncodeToString(sum[:]),
		ExpiresAt:      time.Now().Add(time.Hour).UTC(),
	}
	if err := taskRepo.Create(ctx, task); err != nil {
		t.Fatalf("create upload task: %v", err)
	}

	chunkPath := filepath.Join(tmpRoot, task.UploadID, "chunk_000000.part")
	if err := os.MkdirAll(filepath.Dir(chunkPath), 0o755); err != nil {
		t.Fatalf("create chunk dir: %v", err)
	}
	if err := os.WriteFile(chunkPath, imageBytes, 0o644); err != nil {
		t.Fatalf("write chunk: %v", err)
	}
	if err := chunkRepo.Save(ctx, &model.UploadChunk{
		UploadID:    task.UploadID,
		ChunkIndex:  0,
		Size:        int64(len(imageBytes)),
		Sha256:      task.Sha256,
		StoragePath: filepath.ToSlash(filepath.Join(task.UploadID, "chunk_000000.part")),
	}); err != nil {
		t.Fatalf("save chunk: %v", err)
	}

	svc := NewUploadTaskService(&configs.UploadConfig{}, taskRepo, chunkRepo)
	svc.ConfigureCompletion(NewMergeService(tmpRoot, finalRoot, "/static/uploads/images"), fileRepo)

	record, err := svc.CompleteUpload(ctx, task.UploadID, task.UserID)
	if err != nil {
		t.Fatalf("complete upload: %v", err)
	}
	if record.URL == "" || record.Path == "" {
		t.Fatalf("expected final file paths, got %+v", record)
	}
	if record.ThumbnailURL == "" {
		t.Fatalf("expected thumbnail url to be stored, got %+v", record)
	}
	if _, err := os.Stat(filepath.Join(finalRoot, filepath.FromSlash(record.Path))); err != nil {
		t.Fatalf("expected final file to exist: %v", err)
	}
	thumbnailPath := strings.TrimPrefix(record.ThumbnailURL, "/static/uploads/images/")
	if _, err := os.Stat(filepath.Join(finalRoot, filepath.FromSlash(thumbnailPath))); err != nil {
		t.Fatalf("expected thumbnail file to exist: %v", err)
	}

	reloaded, err := taskRepo.GetByUploadID(ctx, task.UploadID, task.UserID)
	if err != nil {
		t.Fatalf("reload upload task: %v", err)
	}
	if reloaded.Status != model.UploadTaskStatusCompleted {
		t.Fatalf("expected completed status, got %q", reloaded.Status)
	}

	stored, err := fileRepo.GetByID(ctx, record.ID)
	if err != nil {
		t.Fatalf("reload file record: %v", err)
	}
	if stored.ThumbnailURL != record.ThumbnailURL {
		t.Fatalf("expected stored thumbnail url %q, got %q", record.ThumbnailURL, stored.ThumbnailURL)
	}
	if stored.ContentType != "image/png" {
		t.Fatalf("expected stored content type image/png, got %q", stored.ContentType)
	}
}

func TestCompleteUploadDoesNotPersistFileRecordWhenTaskStateChanges(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	ctx := context.Background()
	tmpRoot := t.TempDir()
	finalRoot := t.TempDir()
	taskRepo := db.NewUploadTaskRepository(database)
	chunkRepo := db.NewUploadChunkRepository(database)
	fileRepo := db.NewFileRepository(database)

	body := testPNGBytes(t)
	sum := sha256.Sum256(body)
	task := &model.UploadTask{
		UploadID:       "complete-race",
		UserID:         10,
		BizType:        "article",
		BizID:          "race-1",
		FileName:       "note.png",
		FileSize:       int64(len(body)),
		ContentType:    "image/png",
		ChunkSize:      int64(len(body)),
		ChunkCount:     1,
		UploadedChunks: "0",
		Status:         model.UploadTaskStatusUploading,
		Sha256:         hex.EncodeToString(sum[:]),
		ExpiresAt:      time.Now().Add(time.Hour).UTC(),
	}
	if err := taskRepo.Create(ctx, task); err != nil {
		t.Fatalf("create upload task: %v", err)
	}
	if err := taskRepo.UpdateProgressGuarded(ctx, task.UploadID, task.UserID, task.UploadedChunks, model.UploadTaskStatusCancelled, model.UploadTaskStatusUploading, task.Version); err != nil {
		t.Fatalf("cancel upload task before completion: %v", err)
	}

	chunkPath := filepath.Join(tmpRoot, task.UploadID, "chunk_000000.part")
	if err := os.MkdirAll(filepath.Dir(chunkPath), 0o755); err != nil {
		t.Fatalf("create chunk dir: %v", err)
	}
	if err := os.WriteFile(chunkPath, body, 0o644); err != nil {
		t.Fatalf("write chunk: %v", err)
	}
	if err := chunkRepo.Save(ctx, &model.UploadChunk{UploadID: task.UploadID, ChunkIndex: 0, Size: int64(len(body)), Sha256: task.Sha256, StoragePath: filepath.ToSlash(filepath.Join(task.UploadID, "chunk_000000.part"))}); err != nil {
		t.Fatalf("save chunk: %v", err)
	}

	svc := NewUploadTaskService(&configs.UploadConfig{}, staleTaskStore{task: task}, chunkRepo)
	svc.ConfigureCompletion(NewMergeService(tmpRoot, finalRoot, "/static/uploads/images"), fileRepo)

	if _, err := svc.CompleteUpload(ctx, task.UploadID, task.UserID); !errors.Is(err, db.ErrUploadTaskStateConflict) {
		t.Fatalf("expected guarded completion conflict, got %v", err)
	}
	var count int64
	if err := database.WithContext(ctx).Table("files").Where("upload_id = ?", task.UploadID).Count(&count).Error; err != nil {
		t.Fatalf("count file records: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no durable file record after completion conflict, got %d", count)
	}
	entries, err := os.ReadDir(finalRoot)
	if err != nil {
		t.Fatalf("read final root: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected merged final files to be removed after completion conflict, got %d entries", len(entries))
	}
}

func TestCompleteUploadRejectsMergedContentTypeMismatch(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	ctx := context.Background()
	tmpRoot := t.TempDir()
	finalRoot := t.TempDir()
	taskRepo := db.NewUploadTaskRepository(database)
	chunkRepo := db.NewUploadChunkRepository(database)
	fileRepo := db.NewFileRepository(database)

	body := []byte("not a png")
	sum := sha256.Sum256(body)
	task := &model.UploadTask{
		UploadID:       "complete-type-mismatch",
		UserID:         14,
		FileName:       "cover.png",
		FileSize:       int64(len(body)),
		ContentType:    "image/png",
		ChunkSize:      int64(len(body)),
		ChunkCount:     1,
		UploadedChunks: "0",
		Status:         model.UploadTaskStatusUploading,
		Sha256:         hex.EncodeToString(sum[:]),
		ExpiresAt:      time.Now().Add(time.Hour).UTC(),
	}
	if err := taskRepo.Create(ctx, task); err != nil {
		t.Fatalf("create upload task: %v", err)
	}
	chunkPath := filepath.Join(tmpRoot, task.UploadID, "chunk_000000.part")
	if err := os.MkdirAll(filepath.Dir(chunkPath), 0o755); err != nil {
		t.Fatalf("create chunk dir: %v", err)
	}
	if err := os.WriteFile(chunkPath, body, 0o644); err != nil {
		t.Fatalf("write chunk: %v", err)
	}
	if err := chunkRepo.Save(ctx, &model.UploadChunk{UploadID: task.UploadID, ChunkIndex: 0, Size: int64(len(body)), Sha256: task.Sha256, StoragePath: filepath.ToSlash(filepath.Join(task.UploadID, "chunk_000000.part"))}); err != nil {
		t.Fatalf("save chunk: %v", err)
	}

	svc := NewUploadTaskService(&configs.UploadConfig{}, taskRepo, chunkRepo)
	svc.ConfigureCompletion(NewMergeService(tmpRoot, finalRoot, "/static/uploads/images"), fileRepo)

	if _, err := svc.CompleteUpload(ctx, task.UploadID, task.UserID); err == nil || !strings.Contains(err.Error(), "declared image type") {
		t.Fatalf("expected content type mismatch error, got %v", err)
	}
	var count int64
	if err := database.WithContext(ctx).Table("files").Where("upload_id = ?", task.UploadID).Count(&count).Error; err != nil {
		t.Fatalf("count file records: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no file record after type mismatch, got %d", count)
	}
}

type staleTaskStore struct {
	task *model.UploadTask
}

func (s staleTaskStore) Create(ctx context.Context, task *model.UploadTask) error {
	return nil
}

func (s staleTaskStore) GetByUploadID(ctx context.Context, uploadID string, userID int64) (*model.UploadTask, error) {
	if s.task != nil && s.task.UploadID == uploadID && s.task.UserID == userID {
		return s.task, nil
	}
	return nil, errors.New("task not found")
}

func (s staleTaskStore) UpdateProgressGuarded(ctx context.Context, uploadID string, userID int64, uploadedChunks string, status string, expectedStatus string, expectedVersion int64) error {
	return db.ErrUploadTaskStateConflict
}
func testPNGBytes(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 20), G: uint8(y * 20), B: 180, A: 255})
		}
	}
	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	return out.Bytes()
}
