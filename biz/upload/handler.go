package upload

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	uploadmodel "github.com/Loe1210/personal-site/biz/model/upload"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
	uploadservice "github.com/Loe1210/personal-site/service"
)

func UploadImage(ctx context.Context, c *app.RequestContext) {
	var req uploadmodel.UploadImageRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	header, err := c.FormFile("file")
	if err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := uploadservice.UploadImage(ctx, &req, header)
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

func GetUploadInfo(ctx context.Context, c *app.RequestContext) {
	var req uploadmodel.GetUploadInfoRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := uploadservice.GetUploadInfo(ctx, &req)
	if err != nil {
		if appErr, ok := err.(*errno.AppError); ok {
			response.WriteError(c, appErr)
			return
		}
		response.WriteError(c, errno.Internal)
		return
	}
	if resp == nil || resp.Upload == nil {
		response.WriteError(c, errno.NotFound)
		return
	}

	response.WriteSuccess(c, resp)
}
