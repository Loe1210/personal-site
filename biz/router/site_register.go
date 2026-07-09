package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	siterouter "github.com/Loe1210/personal-site/biz/router/site"
)

func registerSite(h *server.Hertz) {
	siterouter.Register(h)
}