package article

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

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
// @Failure 409 {object} response.Body
// @Security BearerAuth
// @Router /api/admin/articles [post]
func CreateArticle(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.CreateArticleRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := articleservice.CreateArticle(ctx, &req)
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
