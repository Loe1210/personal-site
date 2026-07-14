package service

import (
	"context"
	"testing"
)

type fakeArticleRepo struct{}

func (f *fakeArticleRepo) GetByID(ctx context.Context, id int64) (*ArticleDetail, error) {
	return &ArticleDetail{ID: id, Title: "Hello"}, nil
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
