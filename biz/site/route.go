package site

import "github.com/cloudwego/hertz/pkg/app/server"

func Register(h *server.Hertz) {
	h.GET("/", Home)
	h.GET("/blog", Blog)
	h.GET("/blog/:slug", ArticleDetail)
	h.GET("/about", About)
	h.GET("/admin/login", AdminLogin)
	h.GET("/admin", AdminDashboard)
	h.GET("/admin/articles", AdminArticles)
	h.GET("/admin/articles/new", AdminArticleNew)
	h.GET("/admin/articles/:id/edit", AdminArticleEdit)
	h.GET("/admin/taxonomy", AdminTaxonomy)
	h.GET("/swagger", Swagger)
}
