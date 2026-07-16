package service

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/Loe1210/personal-site/services/media-service/internal/model"
)

type Storage interface {
	Save(name string, content []byte) (string, error)
}

type ChunkStorage interface {
	SaveChunk(uploadID string, chunkIndex int, content io.Reader) (storagePath string, size int64, sha256 string, err error)
	BackupChunk(storagePath string) (backupPath string, exists bool, err error)
	RestoreChunk(storagePath string, backupPath string) error
	DiscardChunkBackup(backupPath string) error
	RemoveChunk(storagePath string) error
}

type Repository interface {
	Save(ctx context.Context, record *model.FileRecord) error
	GetByID(ctx context.Context, id int64) (*model.FileRecord, error)
}

type Service struct {
	storage Storage
	repo    Repository
}

func NewMediaService(storage Storage, repo Repository) *Service {
	return &Service{storage: storage, repo: repo}
}

func (s *Service) Upload(ctx context.Context, in model.UploadInput) (*model.FileRecord, error) {
	if s == nil || s.storage == nil {
		return nil, errors.New("storage is required")
	}
	if strings.TrimSpace(in.FileName) == "" {
		return nil, errors.New("file name is required")
	}
	if len(in.Content) == 0 {
		return nil, errors.New("file content is required")
	}
	if !isAllowedImageContentType(in.ContentType) {
		return nil, errors.New("only image uploads are allowed")
	}
	url, err := s.storage.Save(in.FileName, in.Content)
	if err != nil {
		return nil, err
	}
	record := &model.FileRecord{
		OriginalName: in.FileName,
		URL:          url,
		Path:         url,
		ContentType:  in.ContentType,
		Size:         int64(len(in.Content)),
		BizType:      normalizeBizType(in.BizType),
	}
	if s.repo != nil {
		if err := s.repo.Save(ctx, record); err != nil {
			return nil, err
		}
	}
	return record, nil
}

func (s *Service) GetFile(ctx context.Context, id int64) (*model.FileRecord, error) {
	if s == nil || s.repo == nil {
		return nil, errors.New("repository is required")
	}
	return s.repo.GetByID(ctx, id)
}

func normalizeBizType(input string) string {
	bizType := strings.ToLower(strings.TrimSpace(input))
	if bizType == "" {
		return "common"
	}
	return bizType
}

func isAllowedImageContentType(contentType string) bool {
	switch strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0])) {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml":
		return true
	default:
		return false
	}
}
