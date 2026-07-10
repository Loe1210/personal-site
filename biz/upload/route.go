package upload

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	sessionmw "github.com/Loe1210/personal-site/pkg/middleware/session"
)

func Register(h *server.Hertz) {
	admin := h.Group("/api/admin", sessionmw.AuthMiddleware())
	{
		admin.POST("/upload", UploadImage)
		admin.GET("/uploads/:id", GetUploadInfo)
	}
}
