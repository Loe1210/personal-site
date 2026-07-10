package auth

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	sessionmw "github.com/Loe1210/personal-site/pkg/middleware/session"
)

func Register(h *server.Hertz) {
	admin := h.Group("/api/admin")
	{
		admin.POST("/login", Login)
		admin.POST("/logout", sessionmw.AuthMiddleware(), Logout)
		admin.GET("/me", sessionmw.AuthMiddleware(), Me)
	}
}
