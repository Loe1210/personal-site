package content

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	contentclient "github.com/Loe1210/personal-site/services/gateway/internal/client/content"
)

type Handler struct {
	articles contentclient.ArticleClient
}

func NewHandler(articles contentclient.ArticleClient) *Handler {
	return &Handler{articles: articles}
}

func (h *Handler) ListArticles(ctx context.Context, c *app.RequestContext) {
	resp, err := h.articles.ListPublicArticles(ctx, contentclient.ListPublicArticlesRequest{
		Page:     parseInt64(c.Query("page")),
		PageSize: parseInt64(c.Query("page_size")),
		Keyword:  c.Query("keyword"),
	})
	if err != nil {
		c.JSON(consts.StatusBadGateway, map[string]any{"code": 40003, "message": "content service unavailable"})
		return
	}
	c.JSON(consts.StatusOK, resp)
}

func (h *Handler) GetArticle(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 40004, "message": "invalid article id"})
		return
	}

	article, err := h.articles.GetArticleByID(ctx, id)
	if err != nil {
		c.JSON(consts.StatusBadGateway, map[string]any{"code": 40003, "message": "content service unavailable"})
		return
	}
	c.JSON(consts.StatusOK, article)
}

func parseInt64(value string) int64 {
	parsed, _ := strconv.ParseInt(value, 10, 64)
	return parsed
}
