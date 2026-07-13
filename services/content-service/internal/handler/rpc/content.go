package rpc

import (
	"context"

	"github.com/Loe1210/personal-site/services/content-service/internal/application"
)

type Handler struct {
	articles *application.ArticleService
}

func NewHandler(articles *application.ArticleService) *Handler {
	return &Handler{articles: articles}
}

func (h *Handler) GetArticleByID(ctx context.Context, id int64) (*application.ArticleDetail, error) {
	return h.articles.GetArticleByID(ctx, id)
}

func (h *Handler) ListPublicArticles(ctx context.Context, filter application.ListFilter) (*application.ListResult, error) {
	return h.articles.ListPublicArticles(ctx, filter)
}
