package biz

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/services/auth-service/biz/authenticator"
	"github.com/Loe1210/personal-site/services/auth-service/internal/service"
)

func RegisterRoutes(hertz *server.Hertz, authService *service.Service) {
	handler := authenticator.NewHandler(authService)
	hertz.POST("/login", handler.Login)
	hertz.POST("/logout", handler.Logout)
	hertz.GET("/me", handler.Me)
}
