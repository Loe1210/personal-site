package upload

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/services/media-service/internal/service"
)

func RegisterRoutes(hertz *server.Hertz, media *service.Service, chunks *service.ChunkService) {
	handler := NewHandler(media, chunks)
	hertz.POST("/upload", handler.Upload)
	hertz.GET("/files/:id", handler.GetFile)
	hertz.POST("/upload/tasks/:upload_id/chunks/:chunk_index", handler.UploadChunk)
}
