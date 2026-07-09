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

// UpdateArticle godoc
// @Summary 更新文章
// @Description 根据文章 ID 更新文章内容
// @Tags article-admin
// @Accept json
// @Produce json
// @Param id path int true "文章 ID"
// @Param body body article.UpdateArticleRequest true "更新文章请求"
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 401 {object} response.Body
// @Failure 404 {object} response.Body
// @Security BearerAuth
// @Router /api/admin/articles/{id} [put]
func UpdateArticle(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.UpdateArticleRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	resp, err := articleservice.UpdateArticle(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}
	if resp == nil || resp.Article == nil {
		c.JSON(consts.StatusNotFound, response.Error(errno.ErrorCode, "article not found"))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}