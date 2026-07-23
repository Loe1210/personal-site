package rpc

import (
	"context"

	"github.com/Loe1210/personal-site/internal/xerrors"
	kitexauth "github.com/Loe1210/personal-site/kitex_gen/auth"
	kitexbase "github.com/Loe1210/personal-site/kitex_gen/base"
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
		return &kitexauth.AuthContext{BaseResp: baseResp(xerrors.CodeAuthSessionExpired, "login expired")}, nil
	}
	return &kitexauth.AuthContext{
		UserId:   result.UserID,
		Username: result.Username,
		Roles:    append([]string(nil), result.Roles...),
		BaseResp: baseResp(xerrors.CodeOK, "success"),
	}, nil
}

func (h *Handler) CheckPermission(ctx context.Context, req *kitexauth.CheckPermissionRequest) (*kitexauth.CheckPermissionResponse, error) {
	allowed, err := h.service.CheckPermission(ctx, req.GetUserId(), req.GetCode())
	if err != nil {
		return &kitexauth.CheckPermissionResponse{BaseResp: baseResp(xerrors.CodeAuthPermissionDenied, "permission denied")}, nil
	}
	if !allowed {
		return &kitexauth.CheckPermissionResponse{Allowed: false, BaseResp: baseResp(xerrors.CodeAuthPermissionDenied, "permission denied")}, nil
	}
	return &kitexauth.CheckPermissionResponse{Allowed: allowed, BaseResp: baseResp(xerrors.CodeOK, "success")}, nil
}

func baseResp(code int32, msg string) *kitexbase.BaseResp {
	return &kitexbase.BaseResp{Code: code, Msg: msg}
}
