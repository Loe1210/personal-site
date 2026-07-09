package article

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	articlehandler "github.com/Loe1210/personal-site/biz/handler/article"
	"github.com/Loe1210/personal-site/biz/mw/jwt"
)

func Register(h *server.Hertz) {
	h.GET("/api/articles", articlehandler.ListArticles)
	h.GET("/api/articles/:slug", articlehandler.GetArticleBySlug)

	admin := h.Group("/api/admin", jwt.AuthMiddleware())
	{
		admin.GET("/articles", articlehandler.ListAdminArticles)
		admin.POST("/articles", articlehandler.CreateArticle)
		admin.PUT("/articles/:id", articlehandler.UpdateArticle)
		admin.DELETE("/articles/:id", articlehandler.DeleteArticle)
	}
}