package router

import "github.com/cloudwego/hertz/pkg/app/server"

func Register(h *server.Hertz) {
	registerSite(h)
	registerHealth(h)
	registerAuth(h)
	registerArticle(h)
}
