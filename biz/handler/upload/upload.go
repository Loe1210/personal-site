package upload

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	uploadmodel "github.com/Loe1210/personal-site/biz/model/upload"
	uploadservice "github.com/Loe1210/personal-site/biz/service/upload"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
)

// UploadImage godoc
// @Summary 上传图片
// @Description 后台上传图片文件，返回上传记录和访问 URL
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param biz_type formData string false "业务类型，如 article_cover / site_asset"
// @Param file formData file true "上传图片文件"
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 401 {object} response.Body
// @Failure 500 {object} response.Body
// @Router /api/admin/upload [post]
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

// GetUploadInfo godoc
// @Summary 获取上传文件信息
// @Description 根据上传记录 ID 获取文件元信息
// @Tags upload
// @Accept json
// @Produce json
// @Param id path int true "上传记录 ID"
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 401 {object} response.Body
// @Failure 404 {object} response.Body
// @Router /api/admin/uploads/{id} [get]
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