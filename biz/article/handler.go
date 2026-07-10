package article

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	articlemodel "github.com/Loe1210/personal-site/biz/model/article"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
	articleservice "github.com/Loe1210/personal-site/service"
)

func ListArticles(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.ListArticlesRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	resp, err := articleservice.ListPublicArticles(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}

func GetArticleBySlug(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.GetArticleBySlugRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	resp, err := articleservice.GetPublicArticleBySlug(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}
	if resp == nil || resp.Article == nil {
		c.JSON(consts.StatusNotFound, response.Error(errno.ErrorCode, "article not found"))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}

func ListAdminArticles(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.ListArticlesRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	resp, err := articleservice.ListAdminArticles(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
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
	if resp == nil || resp.Article == nil {
		response.WriteError(c, errno.ArticleNotFound)
		return
	}

	response.WriteSuccess(c, resp)
}

func DeleteArticle(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.DeleteArticleRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	resp, err := articleservice.DeleteArticle(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}
	if resp == nil || !resp.Success {
		c.JSON(consts.StatusNotFound, response.Error(errno.ErrorCode, "article not found"))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}
