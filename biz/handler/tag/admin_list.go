package tag

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	tagmodel "github.com/Loe1210/personal-site/biz/model/tag"
	tagservice "github.com/Loe1210/personal-site/biz/service/tag"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
)

// ListAdminTags godoc
// @Summary 获取后台标签列表
// @Description 获取后台标签列表
// @Tags tag-admin
// @Accept json
// @Produce json
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 401 {object} response.Body
// @Security BearerAuth
// @Router /api/admin/tags [get]
func ListAdminTags(ctx context.Context, c *app.RequestContext) {
	var req tagmodel.ListTagsRequest

	resp, err := tagservice.ListTags(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}