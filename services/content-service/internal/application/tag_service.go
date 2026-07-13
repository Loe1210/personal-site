package application

import "context"

type Tag struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type TagRepository interface {
	List(ctx context.Context) ([]*Tag, error)
	Create(ctx context.Context, tag *Tag) error
	Update(ctx context.Context, tag *Tag) error
	Delete(ctx context.Context, id int64) error
}

type TagService struct {
	repo TagRepository
}

func NewTagService(repo TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) ListTags(ctx context.Context) ([]*Tag, error) {
	return s.repo.List(ctx)
}

func (s *TagService) CreateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	if err := s.repo.Create(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagService) UpdateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	if err := s.repo.Update(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagService) DeleteTag(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
