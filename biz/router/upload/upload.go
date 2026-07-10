package upload

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	uploadhandler "github.com/Loe1210/personal-site/biz/handler/upload"
	authmw "github.com/Loe1210/personal-site/biz/mw/session"
)

func Register(h *server.Hertz) {
	admin := h.Group("/api/admin", authmw.AuthMiddleware())
	{
		admin.POST("/upload", uploadhandler.UploadImage)
		admin.GET("/uploads/:id", uploadhandler.GetUploadInfo)
	}
}