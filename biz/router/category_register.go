package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	categoryrouter "github.com/Loe1210/personal-site/biz/router/category"
)

func registerCategory(h *server.Hertz) {
	categoryrouter.Register(h)
}