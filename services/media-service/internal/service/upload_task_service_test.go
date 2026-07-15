package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
	"github.com/Loe1210/personal-site/services/media-service/internal/model"
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
