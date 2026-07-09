package auth

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	authhandler "github.com/Loe1210/personal-site/biz/handler/auth"
	authmw "github.com/Loe1210/personal-site/biz/mw/jwt"
)

func Register(h *server.Hertz) {
	admin := h.Group("/api/admin")
	{
		admin.POST("/login", authhandler.Login)
		admin.GET("/me", authmw.AuthMiddleware(), authhandler.Me)
	}
}