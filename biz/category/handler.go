package category

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	categorymodel "github.com/Loe1210/personal-site/biz/model/category"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
	categoryservice "github.com/Loe1210/personal-site/service"
)

func ListCategories(ctx context.Context, c *app.RequestContext) {
	var req categorymodel.ListCategoriesRequest

	resp, err := categoryservice.ListCategories(ctx, &req)
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

func ListAdminCategories(ctx context.Context, c *app.RequestContext) {
	var req categorymodel.ListCategoriesRequest

	resp, err := categoryservice.ListCategories(ctx, &req)
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

func CreateCategory(ctx context.Context, c *app.RequestContext) {
	var req categorymodel.CreateCategoryRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := categoryservice.CreateCategory(ctx, &req)
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

func UpdateCategory(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteError(c, errno.BadRequest)
		return
	}

	var req categorymodel.UpdateCategoryRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}
	req.ID = id

	resp, err := categoryservice.UpdateCategory(ctx, &req)
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

func DeleteCategory(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := categoryservice.DeleteCategory(ctx, &categorymodel.DeleteCategoryRequest{ID: id})
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
