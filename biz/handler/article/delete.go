package article

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	articlemodel "github.com/Loe1210/personal-site/biz/model/article"
	articleservice "github.com/Loe1210/personal-site/biz/service/article"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
)

// DeleteArticle godoc
// @Summary 删除文章
// @Description 根据文章 ID 删除文章
// @Tags article-admin
// @Accept json
// @Produce json
// @Param id path int true "文章 ID"
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 401 {object} response.Body
// @Failure 404 {object} response.Body
// @Security BearerAuth
// @Router /api/admin/articles/{id} [delete]
func DeleteArticle(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.DeleteArticleRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	resp, err := articleservice.DeleteArticle(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}
	if resp == nil || !resp.Success {
		c.JSON(consts.StatusNotFound, response.Error(errno.ErrorCode, "article not found"))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}