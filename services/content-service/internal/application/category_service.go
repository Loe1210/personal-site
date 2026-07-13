package application

import "context"

type Category struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type CategoryRepository interface {
	List(ctx context.Context) ([]*Category, error)
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id int64) error
}

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]*Category, error) {
	return s.repo.List(ctx)
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *Category) (*Category, error) {
	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, category *Category) (*Category, error) {
	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
