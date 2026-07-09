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

// GetArticleBySlug godoc
// @Summary 获取文章详情
// @Description 通过 slug 获取已发布文章详情
// @Tags article
// @Accept json
// @Produce json
// @Param slug path string true "文章 slug"
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 404 {object} response.Body
// @Router /api/articles/{slug} [get]
func GetArticleBySlug(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.GetArticleBySlugRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	resp, err := articleservice.GetPublicArticleBySlug(ctx, &req)
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