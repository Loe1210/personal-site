package http

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/Loe1210/personal-site/services/content-service/internal/service"
)

type Handler struct {
	articles *service.ArticleService
}

func NewHandler(articles *service.ArticleService) *Handler {
	return &Handler{articles: articles}
}

func (h *Handler) GetArticleByID(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 10001, "message": "invalid article id"})
		return
	}
	article, err := h.articles.GetArticleByID(ctx, id)
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]any{"code": 30001, "message": "article not found"})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": article})
}

func (h *Handler) ListPublicArticles(ctx context.Context, c *app.RequestContext) {
	result, err := h.articles.ListPublicArticles(ctx, listFilterFromRequest(c))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]any{"code": 50000, "message": "list articles failed"})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": result})
}

func (h *Handler) ListAdminArticles(ctx context.Context, c *app.RequestContext) {
	result, err := h.articles.ListAdminArticles(ctx, listFilterFromRequest(c))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]any{"code": 50000, "message": "list articles failed"})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": result})
}

func (h *Handler) CreateArticle(ctx context.Context, c *app.RequestContext) {
	var article service.ArticleDetail
	if err := c.BindAndValidate(&article); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 10001, "message": "invalid request"})
		return
	}
	created, err := h.articles.CreateArticle(ctx, &article)
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 30002, "message": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": created})
}

func (h *Handler) UpdateArticle(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 10001, "message": "invalid article id"})
		return
	}
	var article service.ArticleDetail
	if err := c.BindAndValidate(&article); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 10001, "message": "invalid request"})
		return
	}
	article.ID = id
	updated, err := h.articles.UpdateArticle(ctx, &article)
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 30003, "message": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": updated})
}

func (h *Handler) DeleteArticle(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 10001, "message": "invalid article id"})
		return
	}
	if err := h.articles.DeleteArticle(ctx, id); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 30004, "message": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success"})
}
