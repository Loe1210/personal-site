package upload

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/services/media-service/internal/service"
)

func RegisterTaskRoutes(hertz *server.Hertz, uploadTasks *service.UploadTaskService) {
	handler := NewTaskHandler(uploadTasks)
	hertz.POST("/upload/tasks/init", handler.InitUpload)
	hertz.GET("/upload/tasks/:upload_id", handler.GetUpload)
	hertz.POST("/upload/tasks/:upload_id/cancel", handler.CancelUpload)
	hertz.POST("/upload/tasks/:upload_id/complete", handler.CompleteUpload)
}
