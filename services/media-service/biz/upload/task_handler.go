package upload

import (
	"context"
	"strconv"

	"github.com/Loe1210/personal-site/services/media-service/internal/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type TaskHandler struct{ service *service.UploadTaskService }

func NewTaskHandler(service *service.UploadTaskService) *TaskHandler {
	return &TaskHandler{service: service}
}
func (h *TaskHandler) InitUpload(ctx context.Context, c *app.RequestContext) {
	userID, err := parseTaskUploadUserID(c)
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20020, "message": err.Error()})
		return
	}
	fileSize, err := parseTaskFormInt64(c, "file_size")
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20021, "message": "invalid file size"})
		return
	}
	chunkSize, err := parseTaskFormInt64(c, "chunk_size")
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20022, "message": "invalid chunk size"})
		return
	}
	task, err := h.service.InitUpload(ctx, service.InitInput{UserID: userID, FileName: c.PostForm("file_name"), FileSize: fileSize, ContentType: c.PostForm("content_type"), BizType: c.PostForm("biz_type"), BizID: c.PostForm("biz_id"), Sha256: c.PostForm("sha256"), ChunkSize: chunkSize})
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20023, "message": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": task})
}
func (h *TaskHandler) GetUpload(ctx context.Context, c *app.RequestContext) {
	userID, err := parseTaskUploadUserID(c)
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20024, "message": err.Error()})
		return
	}
	task, chunks, err := h.service.GetUpload(ctx, c.Param("upload_id"), userID)
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]any{"code": 20025, "message": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": map[string]any{"task": task, "chunks": chunks}})
}
func (h *TaskHandler) CancelUpload(ctx context.Context, c *app.RequestContext) {
	userID, err := parseTaskUploadUserID(c)
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20026, "message": err.Error()})
		return
	}
	if err := h.service.CancelUpload(ctx, c.Param("upload_id"), userID); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20027, "message": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success"})
}
func (h *TaskHandler) CompleteUpload(ctx context.Context, c *app.RequestContext) {
	userID, err := parseTaskUploadUserID(c)
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20028, "message": err.Error()})
		return
	}
	record, err := h.service.CompleteUpload(ctx, c.Param("upload_id"), userID)
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20029, "message": err.Error()})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": record})
}
func parseTaskUploadUserID(c *app.RequestContext) (int64, error) {
	return parseTaskFormInt64(c, "user_id")
}
func parseTaskFormInt64(c *app.RequestContext, key string) (int64, error) {
	value := c.PostForm(key)
	if value == "" {
		value = c.Query(key)
	}
	if value == "" {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseInt(value, 10, 64)
}
