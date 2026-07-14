package rpc

import (
	"context"

	"github.com/Loe1210/personal-site/services/auth-service/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ValidateSession(ctx context.Context, sessionID string) (*service.AuthContext, error) {
	return h.service.ValidateSession(ctx, sessionID)
}

func (h *Handler) CheckPermission(ctx context.Context, userID int64, code string) (bool, error) {
	return h.service.CheckPermission(ctx, userID, code)
}
