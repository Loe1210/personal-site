package service

import (
	"context"
	"errors"

	"github.com/Loe1210/personal-site/services/content-service/internal/model"
)

type Tag = model.Tag

type TagRepository interface {
	ListTags(ctx context.Context) ([]Tag, error)
	CreateTag(ctx context.Context, tag *Tag) error
	UpdateTag(ctx context.Context, tag *Tag) error
	DeleteTag(ctx context.Context, id int64) error
}

type TagService struct {
	repo TagRepository
}

func NewTagService(repo TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) ListTags(ctx context.Context) ([]Tag, error) {
	if s.repo == nil {
		return nil, errors.New("tag repository is required")
	}
	return s.repo.ListTags(ctx)
}

func (s *TagService) CreateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	if s.repo == nil {
		return nil, errors.New("tag repository is required")
	}
	if tag == nil {
		return nil, errors.New("tag is required")
	}
	if err := s.repo.CreateTag(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagService) UpdateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	if s.repo == nil {
		return nil, errors.New("tag repository is required")
	}
	if tag == nil || tag.ID <= 0 {
		return nil, errors.New("tag id is required")
	}
	if err := s.repo.UpdateTag(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagService) DeleteTag(ctx context.Context, id int64) error {
	if s.repo == nil {
		return errors.New("tag repository is required")
	}
	if id <= 0 {
		return errors.New("tag id is required")
	}
	return s.repo.DeleteTag(ctx, id)
}
