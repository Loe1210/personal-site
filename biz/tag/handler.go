package tag

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	tagmodel "github.com/Loe1210/personal-site/biz/model/tag"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
	tagservice "github.com/Loe1210/personal-site/service"
)

func ListTags(ctx context.Context, c *app.RequestContext) {
	var req tagmodel.ListTagsRequest

	resp, err := tagservice.ListTags(ctx, &req)
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

func ListAdminTags(ctx context.Context, c *app.RequestContext) {
	var req tagmodel.ListTagsRequest

	resp, err := tagservice.ListTags(ctx, &req)
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

func CreateTag(ctx context.Context, c *app.RequestContext) {
	var req tagmodel.CreateTagRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := tagservice.CreateTag(ctx, &req)
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

func UpdateTag(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteError(c, errno.BadRequest)
		return
	}

	var req tagmodel.UpdateTagRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}
	req.ID = id

	resp, err := tagservice.UpdateTag(ctx, &req)
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

func DeleteTag(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := tagservice.DeleteTag(ctx, &tagmodel.DeleteTagRequest{ID: id})
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
