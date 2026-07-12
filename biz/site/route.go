package site

import "github.com/cloudwego/hertz/pkg/app/server"

func Register(h *server.Hertz) {
	// 页面路由已在 main.go 中注册，这里仅保留 Swagger
	h.GET("/swagger", Swagger)
}
