package article

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

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
// @Failure 409 {object} response.Body
// @Security BearerAuth
// @Router /api/admin/articles/{id} [put]
func UpdateArticle(ctx context.Context, c *app.RequestContext) {
	var req articlemodel.UpdateArticleRequest
	if err := c.BindAndValidate(&req); err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := articleservice.UpdateArticle(ctx, &req)
	if err != nil {
		if appErr, ok := err.(*errno.AppError); ok {
			response.WriteError(c, appErr)
			return
		}
		response.WriteError(c, errno.Internal)
		return
	}
	if resp == nil || resp.Article == nil {
		response.WriteError(c, errno.ArticleNotFound)
		return
	}

	response.WriteSuccess(c, resp)
}
