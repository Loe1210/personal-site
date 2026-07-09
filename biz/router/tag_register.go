package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	tagrouter "github.com/Loe1210/personal-site/biz/router/tag"
)

func registerTag(h *server.Hertz) {
	tagrouter.Register(h)
}