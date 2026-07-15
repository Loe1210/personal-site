package service

import (
	"context"
	"errors"

	"github.com/Loe1210/personal-site/services/content-service/internal/model"
)

type Category = model.Category

type CategoryRepository interface {
	ListCategories(ctx context.Context) ([]Category, error)
	CreateCategory(ctx context.Context, category *Category) error
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id int64) error
}

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]Category, error) {
	if s.repo == nil {
		return nil, errors.New("category repository is required")
	}
	return s.repo.ListCategories(ctx)
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *Category) (*Category, error) {
	if s.repo == nil {
		return nil, errors.New("category repository is required")
	}
	if category == nil {
		return nil, errors.New("category is required")
	}
	if err := s.repo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, category *Category) (*Category, error) {
	if s.repo == nil {
		return nil, errors.New("category repository is required")
	}
	if category == nil || category.ID <= 0 {
		return nil, errors.New("category id is required")
	}
	if err := s.repo.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int64) error {
	if s.repo == nil {
		return errors.New("category repository is required")
	}
	if id <= 0 {
		return errors.New("category id is required")
	}
	return s.repo.DeleteCategory(ctx, id)
}
