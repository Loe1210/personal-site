package rpc

import (
	"context"
	"time"

	kitexmedia "github.com/Loe1210/personal-site/kitex_gen/media"
	"github.com/Loe1210/personal-site/services/media-service/internal/model"
	"github.com/Loe1210/personal-site/services/media-service/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetFile(ctx context.Context, req *kitexmedia.GetFileRequest) (*kitexmedia.GetFileResponse, error) {
	file, err := h.service.GetFile(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &kitexmedia.GetFileResponse{File: toPBFile(file)}, nil
}

func toPBFile(file *model.FileRecord) *kitexmedia.FileRecord {
	if file == nil {
		return nil
	}
	return &kitexmedia.FileRecord{
		Id:           file.ID,
		OriginalName: file.OriginalName,
		Url:          file.URL,
		Path:         file.Path,
		ContentType:  file.ContentType,
		Size:         file.Size,
		BizType:      file.BizType,
		CreatedAt:    file.CreatedAt.UTC().Format(time.RFC3339),
	}
}
