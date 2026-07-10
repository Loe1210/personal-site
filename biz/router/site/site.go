package site

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	sitehandler "github.com/Loe1210/personal-site/biz/handler/site"
)

func Register(h *server.Hertz) {
	h.GET("/", sitehandler.Home)
	h.GET("/blog", sitehandler.Blog)
	h.GET("/blog/:slug", sitehandler.ArticleDetail)
	h.GET("/about", sitehandler.About)
	h.GET("/swagger", sitehandler.Swagger)
}
