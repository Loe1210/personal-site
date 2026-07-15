package rpc

import (
	"context"

	kitexauth "github.com/Loe1210/personal-site/kitex_gen/auth"
	"github.com/Loe1210/personal-site/services/auth-service/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ValidateSession(ctx context.Context, req *kitexauth.ValidateSessionRequest) (*kitexauth.AuthContext, error) {
	result, err := h.service.ValidateSession(ctx, req.GetSessionId())
	if err != nil {
		return nil, err
	}
	return &kitexauth.AuthContext{
		UserId:   result.UserID,
		Username: result.Username,
		Roles:    append([]string(nil), result.Roles...),
	}, nil
}

func (h *Handler) CheckPermission(ctx context.Context, req *kitexauth.CheckPermissionRequest) (*kitexauth.CheckPermissionResponse, error) {
	allowed, err := h.service.CheckPermission(ctx, req.GetUserId(), req.GetCode())
	if err != nil {
		return nil, err
	}
	return &kitexauth.CheckPermissionResponse{Allowed: allowed}, nil
}
