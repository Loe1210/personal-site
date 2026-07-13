package rpc

import (
	"context"

	"github.com/Loe1210/personal-site/services/auth-service/internal/application"
)

// Handler 是后续 Kitex 服务端的适配层，避免 RPC 类型渗入应用层。
type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ValidateSession(ctx context.Context, sessionID string) (*application.AuthContext, error) {
	return h.service.ValidateSession(ctx, sessionID)
}

func (h *Handler) CheckPermission(ctx context.Context, userID int64, code string) (bool, error) {
	return h.service.CheckPermission(ctx, userID, code)
}
