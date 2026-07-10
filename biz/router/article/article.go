package article

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	articlehandler "github.com/Loe1210/personal-site/biz/handler/article"
	authmw "github.com/Loe1210/personal-site/biz/mw/session"
)

func Register(h *server.Hertz) {
	h.GET("/api/articles", articlehandler.ListArticles)
	h.GET("/api/articles/:slug", articlehandler.GetArticleBySlug)

	admin := h.Group("/api/admin", authmw.AuthMiddleware())
	{
		admin.GET("/articles", authmw.RequirePermission("article:read"), articlehandler.ListAdminArticles)
		admin.POST("/articles", authmw.RequirePermission("article:create"), articlehandler.CreateArticle)
		admin.PUT("/articles/:id", authmw.RequirePermission("article:update"), articlehandler.UpdateArticle)
		admin.DELETE("/articles/:id", authmw.RequirePermission("article:delete"), articlehandler.DeleteArticle)
	}
}