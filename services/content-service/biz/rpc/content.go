package rpc

import (
	"context"

	"github.com/Loe1210/personal-site/services/content-service/internal/service"
)

type Handler struct {
	articles *service.ArticleService
}

func NewHandler(articles *service.ArticleService) *Handler {
	return &Handler{articles: articles}
}

func (h *Handler) GetArticleByID(ctx context.Context, id int64) (*service.ArticleDetail, error) {
	return h.articles.GetArticleByID(ctx, id)
}

func (h *Handler) ListPublicArticles(ctx context.Context, filter service.ListFilter) (*service.ListResult, error) {
	return h.articles.ListPublicArticles(ctx, filter)
}
