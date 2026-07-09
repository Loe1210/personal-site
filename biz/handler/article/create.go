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

// CreateArticle godoc
// @Summary 创建文章
// @Description 创建一篇新文章
// @Tags article-admin
// @Accept json
// @Produce json
// @Param body body article.CreateArticleRequest true "创建文章请求"
// @Success 200 {object} response.Body
// @Failure 400 {object} response.Body
// @Failure 401 {object} response.Body
// @Security BearerAuth
// @Router /api/admin/articles [post]
func CreateArticle(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.CreateArticleRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	resp, err := articleservice.CreateArticle(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrorCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(resp))
}