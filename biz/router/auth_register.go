package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	authrouter "github.com/Loe1210/personal-site/biz/router/auth"
)

func registerAuth(h *server.Hertz) {
	authrouter.Register(h)
}