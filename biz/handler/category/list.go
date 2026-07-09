package category

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	categorymodel "github.com/Loe1210/personal-site/biz/model/category"
	categoryservice "github.com/Loe1210/personal-site/biz/service/category"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
)

// ListCategories godoc
// @Summary 获取分类列表
// @Description 获取公开分类列表
// @Tags category
// @Accept json
// @Produce json
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Router /api/categories [get]
func ListCategories(ctx context.Context, c *app.RequestContext) {
	var req categorymodel.ListCategoriesRequest

	resp, err := categoryservice.ListCategories(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}