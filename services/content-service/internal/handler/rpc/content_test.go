package rpc

import (
	"context"
	"testing"

	kitexcontent "github.com/Loe1210/personal-site/kitex_gen/content"
	"github.com/Loe1210/personal-site/services/content-service/internal/model"
	"github.com/Loe1210/personal-site/services/content-service/internal/service"
)

type fakeArticleRepo struct {
	gotID      int64
	gotFilter  model.ListFilter
	getErr     error
	listErr    error
	getArticle *model.ArticleDetail
	listResult *model.ListResult
}

func (f *fakeArticleRepo) GetByID(_ context.Context, id int64) (*model.ArticleDetail, error) {
	f.gotID = id
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.getArticle, nil
}

func (f *fakeArticleRepo) List(_ context.Context, filter model.ListFilter) (*model.ListResult, error) {
	f.gotFilter = filter
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listResult, nil
}

func (f *fakeArticleRepo) Create(context.Context, *model.ArticleDetail) error {
	return nil
}

func (f *fakeArticleRepo) Update(context.Context, *model.ArticleDetail) error {
	return nil
}

func (f *fakeArticleRepo) Delete(context.Context, int64) error {
	return nil
}

func TestListPublicArticlesMapsRequestAndResponse(t *testing.T) {
	repo := &fakeArticleRepo{
		listResult: &model.ListResult{
			Total: 1,
			List: []*model.ArticleDetail{
				{
					ID:          7,
					Title:       "RPC Article",
					Slug:        "rpc-article",
					Summary:     "summary",
					ContentMd:   "markdown",
					ContentHTML: "<p>html</p>",
					CoverImage:  "/covers/rpc.png",
					CategoryID:  3,
					TagIDs:      []int64{5},
					Status:      "published",
					Tags:        []model.TagDTO{{ID: 5, Name: "Go", Slug: "go"}},
				},
			},
		},
	}
	handler := NewHandler(service.NewArticleService(repo))

	resp, err := handler.ListPublicArticles(context.Background(), &kitexcontent.ListPublicArticlesRequest{
		Page:     2,
		PageSize: 20,
		Keyword:  "go",
	})

	if err != nil {
		t.Fatalf("ListPublicArticles returned error: %v", err)
	}
	if repo.gotFilter.Page != 2 || repo.gotFilter.PageSize != 20 || repo.gotFilter.Keyword != "go" || repo.gotFilter.Status != "published" {
		t.Fatalf("unexpected filter: %#v", repo.gotFilter)
	}
	if resp.GetTotal() != 1 || len(resp.GetList()) != 1 {
		t.Fatalf("unexpected list response: %#v", resp)
	}
	article := resp.GetList()[0]
	if article.GetId() != 7 || article.GetTitle() != "RPC Article" || article.GetContentHtml() != "<p>html</p>" {
		t.Fatalf("unexpected mapped article: %#v", article)
	}
	if len(article.GetTags()) != 1 || article.GetTags()[0].GetSlug() != "go" {
		t.Fatalf("unexpected tags: %#v", article.GetTags())
	}
}

func TestGetArticleByIDMapsRequestAndResponse(t *testing.T) {
	repo := &fakeArticleRepo{
		getArticle: &model.ArticleDetail{ID: 9, Title: "Detail", Slug: "detail", Status: "published"},
	}
	handler := NewHandler(service.NewArticleService(repo))

	resp, err := handler.GetArticleByID(context.Background(), &kitexcontent.GetArticleByIDRequest{Id: 9})

	if err != nil {
		t.Fatalf("GetArticleByID returned error: %v", err)
	}
	if repo.gotID != 9 {
		t.Fatalf("expected id 9, got %d", repo.gotID)
	}
	if resp.GetArticle().GetId() != 9 || resp.GetArticle().GetTitle() != "Detail" {
		t.Fatalf("unexpected article response: %#v", resp.GetArticle())
	}
}

func TestGetArticleByIDReturnsBaseRespWhenMissing(t *testing.T) {
	repo := &fakeArticleRepo{getErr: service.ErrArticleNotFound}
	handler := NewHandler(service.NewArticleService(repo))

	resp, err := handler.GetArticleByID(context.Background(), &kitexcontent.GetArticleByIDRequest{Id: 404})

	if err != nil {
		t.Fatalf("expected nil rpc error, got %v", err)
	}
	if resp.GetBaseResp().GetCode() != 20030001 {
		t.Fatalf("expected article not found base resp, got %#v", resp.GetBaseResp())
	}
}
