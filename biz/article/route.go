package article

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	sessionmw "github.com/Loe1210/personal-site/pkg/middleware/session"
)

func Register(h *server.Hertz) {
	h.GET("/api/articles", ListArticles)
	h.GET("/api/articles/:slug", GetArticleBySlug)

	admin := h.Group("/api/admin", sessionmw.AuthMiddleware())
	{
		admin.GET("/articles", sessionmw.RequirePermission("article:read"), ListAdminArticles)
		admin.POST("/articles", sessionmw.RequirePermission("article:create"), CreateArticle)
		admin.PUT("/articles/:id", sessionmw.RequirePermission("article:update"), UpdateArticle)
		admin.DELETE("/articles/:id", sessionmw.RequirePermission("article:delete"), DeleteArticle)
	}
}
