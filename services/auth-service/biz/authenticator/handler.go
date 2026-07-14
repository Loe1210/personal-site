package authenticator

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/Loe1210/personal-site/pkg/xauth"
	bizmodel "github.com/Loe1210/personal-site/services/auth-service/biz/model"
	"github.com/Loe1210/personal-site/services/auth-service/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(ctx context.Context, c *app.RequestContext) {
	var request bizmodel.LoginRequest
	if err := c.BindAndValidate(&request); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{"code": 10001, "message": "invalid request"})
		return
	}
	bundle, err := h.service.CreateSession(ctx, request.Username, request.Password)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, map[string]any{"code": 10002, "message": "invalid credentials"})
		return
	}
	c.SetCookie(bundle.CookieName, bundle.SessionID, 7200, "/", "", protocol.CookieSameSiteLaxMode, false, true)
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": bundle})
}

func (h *Handler) Logout(_ context.Context, c *app.RequestContext) {
	_ = xauth.DestroySession(xauth.SessionIDFromRequest(c))
	xauth.ClearSessionCookie(c)
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success"})
}

func (h *Handler) Me(ctx context.Context, c *app.RequestContext) {
	user, err := h.service.GetCurrentUser(ctx, xauth.SessionIDFromRequest(c))
	if err != nil {
		c.JSON(consts.StatusUnauthorized, map[string]any{"code": 10002, "message": "login expired"})
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": user})
}
