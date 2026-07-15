package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/Loe1210/personal-site/configs"
	db "github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
	"github.com/Loe1210/personal-site/services/media-service/internal/model"
	"github.com/google/uuid"
)

const defaultChunkSizeBytes int64 = 4 * 1024 * 1024

type InitInput struct {
	UserID      int64
	FileName    string
	FileSize    int64
	ContentType string
	BizType     string
	BizID       string
	Sha256      string
	ChunkSize   int64
}

type UploadTaskRepository interface {
	Create(ctx context.Context, task *model.UploadTask) error
	GetByUploadID(ctx context.Context, uploadID string, userID int64) (*model.UploadTask, error)
	UpdateProgressGuarded(ctx context.Context, uploadID string, userID int64, uploadedChunks string, status string, expectedStatus string, expectedVersion int64) error
}

type UploadTaskService struct {
	cfg         *configs.UploadConfig
	tasks       UploadTaskRepository
	chunks      *db.UploadChunkRepository
	maxUploadSz int64
}

func NewUploadTaskService(cfg *configs.UploadConfig, tasks UploadTaskRepository, chunks *db.UploadChunkRepository) *UploadTaskService {
	svc := &UploadTaskService{cfg: cfg, tasks: tasks, chunks: chunks}
	if cfg != nil && cfg.MaxImageSizeMB > 0 {
		svc.maxUploadSz = cfg.MaxImageSizeMB * 1024 * 1024
	}
	return svc
}

func (s *UploadTaskService) InitUpload(ctx context.Context, in InitInput) (*model.UploadTask, error) {
	if s == nil {
		return nil, errors.New("upload task service is required")
	}
	if in.UserID <= 0 {
		return nil, errors.New("user id is required")
	}
	if in.FileName == "" {
		return nil, errors.New("file name is required")
	}
	if in.FileSize <= 0 {
		return nil, errors.New("file size is required")
	}
	if s.maxUploadSz > 0 && in.FileSize > s.maxUploadSz {
		return nil, errors.New("file too large")
	}
	if s.tasks == nil {
		return nil, errors.New("upload task repository is required")
	}

	chunkSize := in.ChunkSize
	if chunkSize <= 0 {
		chunkSize = defaultChunkSizeBytes
	}
	chunkCount := int(math.Ceil(float64(in.FileSize) / float64(chunkSize)))
	if chunkCount <= 0 {
		chunkCount = 1
	}

	task := &model.UploadTask{
		UploadID:   uuid.NewString(),
		UserID:     in.UserID,
		BizType:    normalizeBizType(in.BizType),
		BizID:      in.BizID,
		FileName:   in.FileName,
		FileSize:   in.FileSize,
		ChunkSize:  chunkSize,
		ChunkCount: chunkCount,
		Status:     model.UploadTaskStatusUploading,
		Sha256:     in.Sha256,
		ExpiresAt:  time.Now().Add(24 * time.Hour).UTC(),
	}
	if err := s.tasks.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *UploadTaskService) GetUpload(ctx context.Context, uploadID string, userID int64) (*model.UploadTask, []model.UploadChunk, error) {
	if s == nil || s.tasks == nil || s.chunks == nil {
		return nil, nil, errors.New("upload task service is not ready")
	}
	task, err := s.tasks.GetByUploadID(ctx, uploadID, userID)
	if err != nil {
		return nil, nil, err
	}
	chunks, err := s.chunks.ListByUploadID(ctx, uploadID)
	if err != nil {
		return nil, nil, err
	}
	return task, chunks, nil
}

func (s *UploadTaskService) CancelUpload(ctx context.Context, uploadID string, userID int64) error {
	return s.updateStatus(ctx, uploadID, userID, model.UploadTaskStatusCancelled)
}

func (s *UploadTaskService) CompleteUpload(ctx context.Context, uploadID string, userID int64) error {
	return s.updateStatus(ctx, uploadID, userID, model.UploadTaskStatusCompleted)
}

func (s *UploadTaskService) updateStatus(ctx context.Context, uploadID string, userID int64, status string) error {
	if s == nil || s.tasks == nil {
		return errors.New("upload task repository is required")
	}
	task, err := s.tasks.GetByUploadID(ctx, uploadID, userID)
	if err != nil {
		return err
	}
	return s.tasks.UpdateProgressGuarded(ctx, task.UploadID, task.UserID, task.UploadedChunks, status, task.Status, task.Version)
}
