package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/services/media-service/biz"
	"github.com/Loe1210/personal-site/services/media-service/internal/service"
)

func newRouter(media *service.Service, addr string) *server.Hertz {
	h := server.Default(server.WithHostPorts(addr))
	biz.RegisterRoutes(h, media)
	return h
}