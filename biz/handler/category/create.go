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

// CreateCategory godoc
// @Summary 创建分类
// @Description 创建一个新分类
// @Tags category-admin
// @Accept json
// @Produce json
// @Param body body category.CreateCategoryRequest true "创建分类请求"
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 401 {object} response.Body
// @Security BearerAuth
// @Router /api/admin/categories [post]
func CreateCategory(ctx context.Context, c *app.RequestContext) {
	var req categorymodel.CreateCategoryRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	resp, err := categoryservice.CreateCategory(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}