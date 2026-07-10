package tag

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	tagmodel "github.com/Loe1210/personal-site/biz/model/tag"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
	tagservice "github.com/Loe1210/personal-site/service"
)

func ListTags(ctx context.Context, c *app.RequestContext) {
	var req tagmodel.ListTagsRequest

	resp, err := tagservice.ListTags(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}

func ListAdminTags(ctx context.Context, c *app.RequestContext) {
	var req tagmodel.ListTagsRequest

	resp, err := tagservice.ListTags(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
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
