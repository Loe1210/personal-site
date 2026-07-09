package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	articlerouter "github.com/Loe1210/personal-site/biz/router/article"
)

func registerArticle(h *server.Hertz) {
	articlerouter.Register(h)
}