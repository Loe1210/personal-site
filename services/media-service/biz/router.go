package biz

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/services/media-service/biz/upload"
	"github.com/Loe1210/personal-site/services/media-service/internal/service"
)

func RegisterRoutes(hertz *server.Hertz, media *service.Service, uploadTasks *service.UploadTaskService) {
	upload.RegisterRoutes(hertz, media)
	upload.RegisterTaskRoutes(hertz, uploadTasks)
}
