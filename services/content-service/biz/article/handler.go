package article

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/internal/xerrors"
	"github.com/Loe1210/personal-site/internal/xhttp"
	bizmodel "github.com/Loe1210/personal-site/services/content-service/biz/model"
	"github.com/Loe1210/personal-site/services/content-service/internal/model"
	"github.com/Loe1210/personal-site/services/content-service/internal/service"
)

type Handler struct {
	articles *service.ArticleService
}

func NewHandler(articles *service.ArticleService) *Handler {
	return &Handler{articles: articles}
}

func RegisterRoutes(hertz *server.Hertz, articles *service.ArticleService) {
	handler := NewHandler(articles)
	hertz.GET("/articles", handler.ListPublicArticles)
	hertz.GET("/articles/:id/adjacent", handler.GetAdjacentArticles)
	hertz.GET("/articles/:id", handler.GetArticleByID)
	hertz.GET("/admin/articles/:id", handler.GetAdminArticleByID)
	hertz.GET("/admin/articles", handler.ListAdminArticles)
	hertz.POST("/admin/articles", handler.CreateArticle)
	hertz.PUT("/admin/articles/:id", handler.UpdateArticle)
	hertz.DELETE("/admin/articles/:id", handler.DeleteArticle)
}

func (h *Handler) GetArticleByID(ctx context.Context, c *app.RequestContext) {
	id, ok := articleIDFromParam(c)
	if !ok {
		return
	}
	article, err := h.articles.GetArticleByID(ctx, id)
	if err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, article)
}

func (h *Handler) GetAdjacentArticles(ctx context.Context, c *app.RequestContext) {
	id, ok := articleIDFromParam(c)
	if !ok {
		return
	}
	adjacent, err := h.articles.GetAdjacentPublicArticles(ctx, id)
	if err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, adjacent)
}

func (h *Handler) GetAdminArticleByID(ctx context.Context, c *app.RequestContext) {
	id, ok := articleIDFromParam(c)
	if !ok {
		return
	}
	article, err := h.articles.GetArticleByID(ctx, id)
	if err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, article)
}

func (h *Handler) ListPublicArticles(ctx context.Context, c *app.RequestContext) {
	result, err := h.articles.ListPublicArticles(ctx, listFilterFromRequest(c))
	if err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, result)
}

func (h *Handler) ListAdminArticles(ctx context.Context, c *app.RequestContext) {
	result, err := h.articles.ListAdminArticles(ctx, listFilterFromRequest(c))
	if err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, result)
}

func (h *Handler) CreateArticle(ctx context.Context, c *app.RequestContext) {
	var article bizmodel.ArticleRequest
	if err := c.BindAndValidate(&article); err != nil {
		xhttp.Fail(c, xerrors.New(xerrors.CodeInvalidArgument, "invalid request"))
		return
	}
	created, err := h.articles.CreateArticle(ctx, toArticleDetail(article))
	if err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, created)
}

func (h *Handler) UpdateArticle(ctx context.Context, c *app.RequestContext) {
	id, ok := articleIDFromParam(c)
	if !ok {
		return
	}
	var article bizmodel.ArticleRequest
	if err := c.BindAndValidate(&article); err != nil {
		xhttp.Fail(c, xerrors.New(xerrors.CodeInvalidArgument, "invalid request"))
		return
	}
	updated, err := h.articles.UpdateArticle(ctx, withID(toArticleDetail(article), id))
	if err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, updated)
}

func (h *Handler) DeleteArticle(ctx context.Context, c *app.RequestContext) {
	id, ok := articleIDFromParam(c)
	if !ok {
		return
	}
	if err := h.articles.DeleteArticle(ctx, id); err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, map[string]bool{"success": true})
}

func articleIDFromParam(c *app.RequestContext) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		xhttp.Fail(c, xerrors.New(xerrors.CodeInvalidArgument, "invalid article id"))
		return 0, false
	}
	return id, true
}

func listFilterFromRequest(c *app.RequestContext) model.ListFilter {
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 64)
	return model.ListFilter{
		Page:     page,
		PageSize: pageSize,
		Status:   c.Query("status"),
		Keyword:  c.Query("keyword"),
	}
}

func toArticleDetail(article bizmodel.ArticleRequest) *model.ArticleDetail {
	return &model.ArticleDetail{
		Title:       article.Title,
		Slug:        article.Slug,
		Summary:     article.Summary,
		ContentMd:   article.ContentMd,
		ContentHTML: article.ContentHTML,
		CoverImage:  article.CoverImage,
		CategoryID:  article.CategoryID,
		TagIDs:      article.TagIDs,
		Status:      article.Status,
	}
}

func withID(article *model.ArticleDetail, id int64) *model.ArticleDetail {
	article.ID = id
	return article
}
