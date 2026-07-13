package http

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/Loe1210/personal-site/services/web-bff/internal/assembler"
)

type Handler struct {
	articles *assembler.ArticlePageAssembler
}

func NewHandler(articles *assembler.ArticlePageAssembler) *Handler {
	return &Handler{articles: articles}
}

func (h *Handler) RegisterRoutes(hertz *server.Hertz) {
	hertz.GET("/blog/articles/:id", h.GetArticlePage)
}

func (h *Handler) GetArticlePage(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 10001, "message": "invalid article id"})
		return
	}
	page, err := h.articles.BuildArticlePage(ctx, id)
	if err != nil {
		c.JSON(consts.StatusBadGateway, map[string]any{"code": 40001, "message": "build article page failed"})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": page})
}
