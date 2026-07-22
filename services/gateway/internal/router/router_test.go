package router

import (
	"context"
	"testing"

	contentclient "github.com/Loe1210/personal-site/services/gateway/internal/client/content"
	contenthandler "github.com/Loe1210/personal-site/services/gateway/internal/handler/content"
)

type fakeArticleClient struct{}

func (f *fakeArticleClient) ListPublicArticles(context.Context, contentclient.ListPublicArticlesRequest) (*contentclient.ListPublicArticlesResponse, error) {
	return &contentclient.ListPublicArticlesResponse{}, nil
}

func (f *fakeArticleClient) GetArticleByID(context.Context, int64) (*contentclient.Article, error) {
	return &contentclient.Article{}, nil
}

func TestRegisterRoutesRequiresDependencies(t *testing.T) {
	deps := Dependencies{
		AuthServiceName: "auth-service",
		ContentHandler:  contenthandler.NewHandler(&fakeArticleClient{}),
	}
	if err := ValidateDependencies(deps); err != nil {
		t.Fatalf("expected dependencies to validate, got %v", err)
	}
}

func TestValidateDependenciesRequiresContentHandler(t *testing.T) {
	deps := Dependencies{AuthServiceName: "auth-service"}
	if err := ValidateDependencies(deps); err == nil {
		t.Fatal("expected missing content handler to fail validation")
	}
}
