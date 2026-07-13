package application

import (
	"context"
	"errors"
	"strings"

	"github.com/Loe1210/personal-site/services/media-service/internal/domain"
)

type Storage interface {
	Save(name string, content []byte) (string, error)
}

type Repository interface {
	Save(ctx context.Context, record *domain.FileRecord) error
	GetByID(ctx context.Context, id int64) (*domain.FileRecord, error)
}

type UploadInput struct {
	FileName    string
	Content     []byte
	ContentType string
	BizType     string
}

type FileRecord = domain.FileRecord

type Service struct {
	storage Storage
	repo    Repository
}

func NewMediaService(storage Storage, repo Repository) *Service {
	return &Service{storage: storage, repo: repo}
}

func (s *Service) Upload(ctx context.Context, in UploadInput) (*FileRecord, error) {
	if s == nil || s.storage == nil {
		return nil, errors.New("storage is required")
	}
	if strings.TrimSpace(in.FileName) == "" {
		return nil, errors.New("file name is required")
	}
	if len(in.Content) == 0 {
		return nil, errors.New("file content is required")
	}
	url, err := s.storage.Save(in.FileName, in.Content)
	if err != nil {
		return nil, err
	}
	record := &domain.FileRecord{
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

func (s *Service) GetFile(ctx context.Context, id int64) (*FileRecord, error) {
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
