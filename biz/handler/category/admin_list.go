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

// ListAdminCategories godoc
// @Summary 获取后台分类列表
// @Description 获取后台分类列表
// @Tags category-admin
// @Accept json
// @Produce json
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 401 {object} response.Body
// @Security BearerAuth
// @Router /api/admin/categories [get]
func ListAdminCategories(ctx context.Context, c *app.RequestContext) {
	var req categorymodel.ListCategoriesRequest

	resp, err := categoryservice.ListCategories(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}