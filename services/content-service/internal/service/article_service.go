package service

import (
	"context"
	"errors"

	"github.com/Loe1210/personal-site/services/content-service/internal/model"
)

type ArticleDetail = model.ArticleDetail
type TagDTO = model.TagDTO
type ListFilter = model.ListFilter
type ListResult = model.ListResult

type ArticleGetter interface {
	GetByID(ctx context.Context, id int64) (*ArticleDetail, error)
}

type ArticleLister interface {
	List(ctx context.Context, filter ListFilter) (*ListResult, error)
}

type ArticleWriter interface {
	Create(ctx context.Context, article *ArticleDetail) error
	Update(ctx context.Context, article *ArticleDetail) error
	Delete(ctx context.Context, id int64) error
}

type ArticleRepository interface {
	ArticleGetter
	ArticleLister
	ArticleWriter
}

type ArticleService struct {
	getter ArticleGetter
	lister ArticleLister
	writer ArticleWriter
}

func NewArticleService(repo ArticleGetter) *ArticleService {
	service := &ArticleService{getter: repo}
	if lister, ok := repo.(ArticleLister); ok {
		service.lister = lister
	}
	if writer, ok := repo.(ArticleWriter); ok {
		service.writer = writer
	}
	return service
}

func (s *ArticleService) GetArticleByID(ctx context.Context, id int64) (*ArticleDetail, error) {
	if id <= 0 {
		return nil, errors.New("article id is required")
	}
	return s.getter.GetByID(ctx, id)
}

func (s *ArticleService) ListPublicArticles(ctx context.Context, filter ListFilter) (*ListResult, error) {
	filter.Status = "published"
	if s.lister == nil {
		return nil, errors.New("article lister is required")
	}
	return s.lister.List(ctx, normalizeListFilter(filter))
}

func (s *ArticleService) ListAdminArticles(ctx context.Context, filter ListFilter) (*ListResult, error) {
	if s.lister == nil {
		return nil, errors.New("article lister is required")
	}
	return s.lister.List(ctx, normalizeListFilter(filter))
}

func (s *ArticleService) CreateArticle(ctx context.Context, article *ArticleDetail) (*ArticleDetail, error) {
	if article == nil {
		return nil, errors.New("article is required")
	}
	if s.writer == nil {
		return nil, errors.New("article writer is required")
	}
	if article.Status == "" {
		article.Status = "draft"
	}
	if err := s.writer.Create(ctx, article); err != nil {
		return nil, err
	}
	return article, nil
}

func (s *ArticleService) UpdateArticle(ctx context.Context, article *ArticleDetail) (*ArticleDetail, error) {
	if article == nil || article.ID <= 0 {
		return nil, errors.New("article id is required")
	}
	if s.writer == nil {
		return nil, errors.New("article writer is required")
	}
	if err := s.writer.Update(ctx, article); err != nil {
		return nil, err
	}
	return s.getter.GetByID(ctx, article.ID)
}

func (s *ArticleService) DeleteArticle(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("article id is required")
	}
	if s.writer == nil {
		return errors.New("article writer is required")
	}
	return s.writer.Delete(ctx, id)
}

func normalizeListFilter(filter ListFilter) ListFilter {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}
	return filter
}
