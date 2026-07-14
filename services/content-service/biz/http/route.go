package http

import "github.com/cloudwego/hertz/pkg/app/server"

func (h *Handler) RegisterRoutes(hertz *server.Hertz) {
	hertz.GET("/articles", h.ListPublicArticles)
	hertz.GET("/articles/:id", h.GetArticleByID)
	hertz.GET("/admin/articles", h.ListAdminArticles)
	hertz.POST("/admin/articles", h.CreateArticle)
	hertz.PUT("/admin/articles/:id", h.UpdateArticle)
	hertz.DELETE("/admin/articles/:id", h.DeleteArticle)
}
