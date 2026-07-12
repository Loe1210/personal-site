package article

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"

	articlemodel "github.com/Loe1210/personal-site/biz/model/article"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
	articleservice "github.com/Loe1210/personal-site/service"
)

func ListArticles(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.ListArticlesRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := articleservice.ListPublicArticles(ctx, &req)
	if err != nil {
		if appErr, ok := err.(*errno.AppError); ok {
			response.WriteError(c, appErr)
			return
		}
		response.WriteError(c, errno.Internal)
		return
	}

	response.WriteSuccess(c, resp)
}

func GetArticleBySlug(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.GetArticleBySlugRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := articleservice.GetPublicArticleBySlug(ctx, &req)
	if err != nil {
		if appErr, ok := err.(*errno.AppError); ok {
			response.WriteError(c, appErr)
			return
		}
		response.WriteError(c, errno.Internal)
		return
	}

	response.WriteSuccess(c, resp)
}

func GetArticleByID(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	var id int64
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil || id <= 0 {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := articleservice.GetPublicArticleByID(ctx, id)
	if err != nil {
		if appErr, ok := err.(*errno.AppError); ok {
			response.WriteError(c, appErr)
			return
		}
		response.WriteError(c, errno.Internal)
		return
	}

	response.WriteSuccess(c, resp)
}

func ListAdminArticles(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.ListArticlesRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := articleservice.ListAdminArticles(ctx, &req)
	if err != nil {
		if appErr, ok := err.(*errno.AppError); ok {
			response.WriteError(c, appErr)
			return
		}
		response.WriteError(c, errno.Internal)
		return
	}

	response.WriteSuccess(c, resp)
}

func CreateArticle(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.CreateArticleRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := articleservice.CreateArticle(ctx, &req)
	if err != nil {
		if appErr, ok := err.(*errno.AppError); ok {
			response.WriteError(c, appErr)
			return
		}
		response.WriteError(c, errno.Internal)
		return
	}

	response.WriteSuccess(c, resp)
}

func UpdateArticle(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.UpdateArticleRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := articleservice.UpdateArticle(ctx, &req)
	if err != nil {
		if appErr, ok := err.(*errno.AppError); ok {
			response.WriteError(c, appErr)
			return
		}
		response.WriteError(c, errno.Internal)
		return
	}

	response.WriteSuccess(c, resp)
}

func DeleteArticle(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.DeleteArticleRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := articleservice.DeleteArticle(ctx, &req)
	if err != nil {
		if appErr, ok := err.(*errno.AppError); ok {
			response.WriteError(c, appErr)
			return
		}
		response.WriteError(c, errno.Internal)
		return
	}

	response.WriteSuccess(c, resp)
}
