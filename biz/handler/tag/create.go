package tag

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	tagmodel "github.com/Loe1210/personal-site/biz/model/tag"
	tagservice "github.com/Loe1210/personal-site/biz/service/tag"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
)

// CreateTag godoc
// @Summary 创建标签
// @Description 创建一个新标签
// @Tags tag-admin
// @Accept json
// @Produce json
// @Param body body tag.CreateTagRequest true "创建标签请求"
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 401 {object} response.Body
// @Failure 409 {object} response.Body
// @Security BearerAuth
// @Router /api/admin/tags [post]
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
