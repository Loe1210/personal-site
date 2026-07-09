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

// ListTags godoc
// @Summary 获取标签列表
// @Description 获取公开标签列表
// @Tags tag
// @Accept json
// @Produce json
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Router /api/tags [get]
func ListTags(ctx context.Context, c *app.RequestContext) {
	var req tagmodel.ListTagsRequest

	resp, err := tagservice.ListTags(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}