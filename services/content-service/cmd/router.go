package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/services/content-service/biz"
	"github.com/Loe1210/personal-site/services/content-service/internal/service"
)

func newRouter(articles *service.ArticleService, addr string) *server.Hertz {
	h := server.Default(server.WithHostPorts(addr))
	biz.RegisterRoutes(h, articles)
	return h
}