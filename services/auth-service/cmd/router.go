package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/services/auth-service/biz"
	"github.com/Loe1210/personal-site/services/auth-service/internal/service"
)

func registerRoutes(hertz *server.Hertz, authService *service.Service) {
	biz.RegisterRoutes(hertz, authService)
}
