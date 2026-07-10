package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	uploadrouter "github.com/Loe1210/personal-site/biz/router/upload"
)

func registerUpload(h *server.Hertz) {
	uploadrouter.Register(h)
}