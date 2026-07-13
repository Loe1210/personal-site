package assembler

import (
	"context"
	"testing"
)

type fakeContentClient struct{}

func (f *fakeContentClient) GetArticleByID(ctx context.Context, id int64) (*ArticleDTO, error) {
	return &ArticleDTO{ID: id, Title: "Hello"}, nil
}

func TestBuildArticlePage(t *testing.T) {
	assembler := NewArticlePageAssembler(&fakeContentClient{})
	page, err := assembler.BuildArticlePage(context.Background(), 12)
	if err != nil {
		t.Fatalf("BuildArticlePage returned error: %v", err)
	}
	if page.Article.ID != 12 {
		t.Fatalf("expected article id 12, got %d", page.Article.ID)
	}
}
