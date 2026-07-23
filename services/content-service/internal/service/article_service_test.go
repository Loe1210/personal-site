package service

import (
	"context"
	"testing"

	"github.com/Loe1210/personal-site/internal/xerrors"
)

type fakeArticleRepo struct {
	calls   int
	article *ArticleDetail
	err     error
}

func (f *fakeArticleRepo) GetByID(ctx context.Context, id int64) (*ArticleDetail, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	if f.article != nil {
		return f.article, nil
	}
	return &ArticleDetail{ID: id, Title: "Hello"}, nil
}

type fakeArticleCache struct {
	items map[int64]*ArticleDetail
	gets  int
	sets  int
	dels  int
}

func newFakeArticleCache() *fakeArticleCache {
	return &fakeArticleCache{items: map[int64]*ArticleDetail{}}
}

func (f *fakeArticleCache) GetArticle(ctx context.Context, id int64) (*ArticleDetail, bool, error) {
	f.gets++
	article, ok := f.items[id]
	return article, ok, nil
}

func (f *fakeArticleCache) SetArticle(ctx context.Context, article *ArticleDetail) error {
	f.sets++
	f.items[article.ID] = article
	return nil
}

func (f *fakeArticleCache) DeleteArticle(ctx context.Context, id int64) error {
	f.dels++
	delete(f.items, id)
	return nil
}

func TestGetArticleByID(t *testing.T) {
	svc := NewArticleService(&fakeArticleRepo{})
	article, err := svc.GetArticleByID(context.Background(), 12)
	if err != nil {
		t.Fatalf("GetArticleByID returned error: %v", err)
	}
	if article.ID != 12 {
		t.Fatalf("expected article id 12, got %d", article.ID)
	}
}

func TestGetArticleByIDReturnsBusinessNotFound(t *testing.T) {
	repo := &fakeArticleRepo{err: ErrArticleNotFound}
	svc := NewArticleService(repo)

	_, err := svc.GetArticleByID(context.Background(), 404)

	if xerrors.CodeOf(err) != xerrors.CodeContentArticleNotFound {
		t.Fatalf("expected article not found code, got %v", err)
	}
}

func TestGetArticleByIDUsesLocalCacheBeforeRedisAndRepository(t *testing.T) {
	local := newFakeArticleCache()
	redis := newFakeArticleCache()
	repo := &fakeArticleRepo{}
	local.items[7] = &ArticleDetail{ID: 7, Title: "local"}
	svc := NewArticleServiceWithCaches(repo, local, redis)

	article, err := svc.GetArticleByID(context.Background(), 7)

	if err != nil {
		t.Fatalf("GetArticleByID returned error: %v", err)
	}
	if article.Title != "local" {
		t.Fatalf("expected local article, got %#v", article)
	}
	if redis.gets != 0 || repo.calls != 0 {
		t.Fatalf("expected local hit only, redis gets=%d repo calls=%d", redis.gets, repo.calls)
	}
}

func TestGetArticleByIDBackfillsLocalFromRedisHit(t *testing.T) {
	local := newFakeArticleCache()
	redis := newFakeArticleCache()
	repo := &fakeArticleRepo{}
	redis.items[8] = &ArticleDetail{ID: 8, Title: "redis"}
	svc := NewArticleServiceWithCaches(repo, local, redis)

	article, err := svc.GetArticleByID(context.Background(), 8)

	if err != nil {
		t.Fatalf("GetArticleByID returned error: %v", err)
	}
	if article.Title != "redis" {
		t.Fatalf("expected redis article, got %#v", article)
	}
	if local.sets != 1 || repo.calls != 0 {
		t.Fatalf("expected local backfill and no repo call, local sets=%d repo calls=%d", local.sets, repo.calls)
	}
}

func TestGetArticleByIDBackfillsCachesFromRepository(t *testing.T) {
	local := newFakeArticleCache()
	redis := newFakeArticleCache()
	repo := &fakeArticleRepo{article: &ArticleDetail{ID: 9, Title: "mysql"}}
	svc := NewArticleServiceWithCaches(repo, local, redis)

	article, err := svc.GetArticleByID(context.Background(), 9)

	if err != nil {
		t.Fatalf("GetArticleByID returned error: %v", err)
	}
	if article.Title != "mysql" {
		t.Fatalf("expected mysql article, got %#v", article)
	}
	if repo.calls != 1 || local.sets != 1 || redis.sets != 1 {
		t.Fatalf("expected mysql backfill, repo calls=%d local sets=%d redis sets=%d", repo.calls, local.sets, redis.sets)
	}
}
