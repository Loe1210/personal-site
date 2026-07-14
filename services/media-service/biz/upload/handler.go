package upload

import (
	"context"
	"io"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/Loe1210/personal-site/services/media-service/internal/model"
	"github.com/Loe1210/personal-site/services/media-service/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Upload(ctx context.Context, c *app.RequestContext) {
	header, err := c.FormFile("file")
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20009, "message": "upload file is required"})
		return
	}
	file, err := header.Open()
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]any{"code": 50000, "message": "open upload file failed"})
		return
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]any{"code": 50000, "message": "read upload file failed"})
		return
	}
	record, err := h.service.Upload(ctx, model.UploadInput{
		FileName:    header.Filename,
		Content:     content,
		ContentType: string(header.Header.Get("Content-Type")),
		BizType:     c.PostForm("biz_type"),
	})
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20010, "message": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": record})
}

func (h *Handler) GetFile(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20013, "message": "invalid file id"})
		return
	}
	record, err := h.service.GetFile(ctx, id)
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]any{"code": 20014, "message": "file not found"})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": record})
}