package meta

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/internal/xerrors"
	"github.com/Loe1210/personal-site/internal/xhttp"
	"github.com/Loe1210/personal-site/services/content-service/internal/model"
	"github.com/Loe1210/personal-site/services/content-service/internal/service"
)

type Handler struct {
	categories *service.CategoryService
	tags       *service.TagService
}

type metaRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

func NewHandler(categories *service.CategoryService, tags *service.TagService) *Handler {
	return &Handler{categories: categories, tags: tags}
}

func RegisterRoutes(hertz *server.Hertz, categories *service.CategoryService, tags *service.TagService) {
	handler := NewHandler(categories, tags)
	hertz.GET("/categories", handler.ListCategories)
	hertz.POST("/categories", handler.CreateCategory)
	hertz.PUT("/categories/:id", handler.UpdateCategory)
	hertz.DELETE("/categories/:id", handler.DeleteCategory)

	hertz.GET("/tags", handler.ListTags)
	hertz.POST("/tags", handler.CreateTag)
	hertz.PUT("/tags/:id", handler.UpdateTag)
	hertz.DELETE("/tags/:id", handler.DeleteTag)

	hertz.GET("/admin/categories", handler.ListCategories)
	hertz.POST("/admin/categories", handler.CreateCategory)
	hertz.PUT("/admin/categories/:id", handler.UpdateCategory)
	hertz.DELETE("/admin/categories/:id", handler.DeleteCategory)

	hertz.GET("/admin/tags", handler.ListTags)
	hertz.POST("/admin/tags", handler.CreateTag)
	hertz.PUT("/admin/tags/:id", handler.UpdateTag)
	hertz.DELETE("/admin/tags/:id", handler.DeleteTag)
}

func (h *Handler) ListCategories(ctx context.Context, c *app.RequestContext) {
	items, err := h.categories.ListCategories(ctx)
	if err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, map[string]any{"list": items})
}

func (h *Handler) CreateCategory(ctx context.Context, c *app.RequestContext) {
	req, ok := bindMetaRequest(c)
	if !ok {
		return
	}
	created, err := h.categories.CreateCategory(ctx, &model.Category{Name: req.Name, Slug: req.Slug, Description: req.Description})
	if err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, created)
}

func (h *Handler) UpdateCategory(ctx context.Context, c *app.RequestContext) {
	id, ok := metaIDFromParam(c, "category")
	if !ok {
		return
	}
	req, ok := bindMetaRequest(c)
	if !ok {
		return
	}
	updated, err := h.categories.UpdateCategory(ctx, &model.Category{ID: id, Name: req.Name, Slug: req.Slug, Description: req.Description})
	if err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, updated)
}

func (h *Handler) DeleteCategory(ctx context.Context, c *app.RequestContext) {
	id, ok := metaIDFromParam(c, "category")
	if !ok {
		return
	}
	if err := h.categories.DeleteCategory(ctx, id); err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, map[string]bool{"success": true})
}

func (h *Handler) ListTags(ctx context.Context, c *app.RequestContext) {
	items, err := h.tags.ListTags(ctx)
	if err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, map[string]any{"list": items})
}

func (h *Handler) CreateTag(ctx context.Context, c *app.RequestContext) {
	req, ok := bindMetaRequest(c)
	if !ok {
		return
	}
	created, err := h.tags.CreateTag(ctx, &model.Tag{Name: req.Name, Slug: req.Slug, Description: req.Description})
	if err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, created)
}

func (h *Handler) UpdateTag(ctx context.Context, c *app.RequestContext) {
	id, ok := metaIDFromParam(c, "tag")
	if !ok {
		return
	}
	req, ok := bindMetaRequest(c)
	if !ok {
		return
	}
	updated, err := h.tags.UpdateTag(ctx, &model.Tag{ID: id, Name: req.Name, Slug: req.Slug, Description: req.Description})
	if err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, updated)
}

func (h *Handler) DeleteTag(ctx context.Context, c *app.RequestContext) {
	id, ok := metaIDFromParam(c, "tag")
	if !ok {
		return
	}
	if err := h.tags.DeleteTag(ctx, id); err != nil {
		xhttp.Fail(c, err)
		return
	}
	xhttp.OK(c, map[string]bool{"success": true})
}

func bindMetaRequest(c *app.RequestContext) (metaRequest, bool) {
	var req metaRequest
	if err := c.BindAndValidate(&req); err != nil || req.Name == "" || req.Slug == "" {
		xhttp.Fail(c, xerrors.New(xerrors.CodeInvalidArgument, "invalid request"))
		return metaRequest{}, false
	}
	return req, true
}

func metaIDFromParam(c *app.RequestContext, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		xhttp.Fail(c, xerrors.New(xerrors.CodeInvalidArgument, "invalid "+name+" id"))
		return 0, false
	}
	return id, true
}
