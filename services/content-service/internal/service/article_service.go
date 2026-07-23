package service

import (
	"context"
	"errors"

	"github.com/Loe1210/personal-site/internal/xerrors"
	"github.com/Loe1210/personal-site/services/content-service/internal/model"
	"gorm.io/gorm"
)

var ErrArticleNotFound = xerrors.New(xerrors.CodeContentArticleNotFound, "article not found")

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

type ArticleAdjacentFinder interface {
	GetAdjacentPublic(ctx context.Context, id int64) (*model.AdjacentArticles, error)
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
	getter      ArticleGetter
	lister      ArticleLister
	adjacent    ArticleAdjacentFinder
	writer      ArticleWriter
	localCache  ArticleCache
	remoteCache ArticleCache
}

func NewArticleService(repo ArticleGetter) *ArticleService {
	return NewArticleServiceWithCaches(repo, nil, nil)
}

func NewArticleServiceWithCaches(repo ArticleGetter, localCache ArticleCache, remoteCache ArticleCache) *ArticleService {
	service := &ArticleService{getter: repo, localCache: localCache, remoteCache: remoteCache}
	if lister, ok := repo.(ArticleLister); ok {
		service.lister = lister
	}
	if adjacent, ok := repo.(ArticleAdjacentFinder); ok {
		service.adjacent = adjacent
	}
	if writer, ok := repo.(ArticleWriter); ok {
		service.writer = writer
	}
	return service
}

func (s *ArticleService) GetArticleByID(ctx context.Context, id int64) (*ArticleDetail, error) {
	if id <= 0 {
		return nil, xerrors.New(xerrors.CodeInvalidArgument, "article id is required")
	}
	if article, ok := s.getFromCache(ctx, s.localCache, id); ok {
		return article, nil
	}
	if article, ok := s.getFromCache(ctx, s.remoteCache, id); ok {
		s.setCache(ctx, s.localCache, article)
		return article, nil
	}
	article, err := s.getter.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrArticleNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	s.setCache(ctx, s.remoteCache, article)
	s.setCache(ctx, s.localCache, article)
	return article, nil
}

func (s *ArticleService) ListPublicArticles(ctx context.Context, filter ListFilter) (*ListResult, error) {
	filter.Status = "published"
	if s.lister == nil {
		return nil, errors.New("article lister is required")
	}
	return s.lister.List(ctx, normalizeListFilter(filter))
}

func (s *ArticleService) GetAdjacentPublicArticles(ctx context.Context, id int64) (*model.AdjacentArticles, error) {
	if id <= 0 {
		return nil, xerrors.New(xerrors.CodeInvalidArgument, "article id is required")
	}
	if s.adjacent == nil {
		return nil, errors.New("article adjacent finder is required")
	}
	result, err := s.adjacent.GetAdjacentPublic(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	return result, nil
}

func (s *ArticleService) ListAdminArticles(ctx context.Context, filter ListFilter) (*ListResult, error) {
	if s.lister == nil {
		return nil, errors.New("article lister is required")
	}
	return s.lister.List(ctx, normalizeListFilter(filter))
}

func (s *ArticleService) CreateArticle(ctx context.Context, article *ArticleDetail) (*ArticleDetail, error) {
	if article == nil {
		return nil, xerrors.New(xerrors.CodeInvalidArgument, "article is required")
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
	s.deleteArticleCaches(ctx, article.ID)
	return article, nil
}

func (s *ArticleService) UpdateArticle(ctx context.Context, article *ArticleDetail) (*ArticleDetail, error) {
	if article == nil || article.ID <= 0 {
		return nil, xerrors.New(xerrors.CodeInvalidArgument, "article id is required")
	}
	if s.writer == nil {
		return nil, errors.New("article writer is required")
	}
	if err := s.writer.Update(ctx, article); err != nil {
		return nil, err
	}
	s.deleteArticleCaches(ctx, article.ID)
	return s.GetArticleByID(ctx, article.ID)
}

func (s *ArticleService) DeleteArticle(ctx context.Context, id int64) error {
	if id <= 0 {
		return xerrors.New(xerrors.CodeInvalidArgument, "article id is required")
	}
	if s.writer == nil {
		return errors.New("article writer is required")
	}
	if err := s.writer.Delete(ctx, id); err != nil {
		return err
	}
	s.deleteArticleCaches(ctx, id)
	return nil
}

func (s *ArticleService) getFromCache(ctx context.Context, cache ArticleCache, id int64) (*ArticleDetail, bool) {
	if cache == nil {
		return nil, false
	}
	article, ok, err := cache.GetArticle(ctx, id)
	if err != nil {
		return nil, false
	}
	return article, ok
}

func (s *ArticleService) setCache(ctx context.Context, cache ArticleCache, article *ArticleDetail) {
	if cache != nil {
		_ = cache.SetArticle(ctx, article)
	}
}

func (s *ArticleService) deleteArticleCaches(ctx context.Context, id int64) {
	if s.localCache != nil {
		_ = s.localCache.DeleteArticle(ctx, id)
	}
	if s.remoteCache != nil {
		_ = s.remoteCache.DeleteArticle(ctx, id)
	}
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
