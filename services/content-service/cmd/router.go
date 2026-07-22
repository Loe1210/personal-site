package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/internal/xhttp"
	"github.com/Loe1210/personal-site/services/content-service/biz"
	"github.com/Loe1210/personal-site/services/content-service/internal/service"
)

func newRouter(articles *service.ArticleService, categories *service.CategoryService, tags *service.TagService, addr string) *server.Hertz {
	h := server.Default(server.WithHostPorts(addr))
	h.Use(xhttp.Recover())
	biz.RegisterRoutes(h, articles, categories, tags)
	return h
}
