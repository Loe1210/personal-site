package authenticator

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"

	"github.com/Loe1210/personal-site/internal/xerrors"
	"github.com/Loe1210/personal-site/internal/xhttp"
	bizmodel "github.com/Loe1210/personal-site/services/auth-service/biz/model"
	"github.com/Loe1210/personal-site/services/auth-service/internal/service"
	"github.com/Loe1210/personal-site/services/auth-service/pkg/xauth"
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
		xhttp.Fail(c, xerrors.New(xerrors.CodeInvalidArgument, "invalid request"))
		return
	}
	bundle, err := h.service.CreateSession(ctx, request.Username, request.Password)
	if err != nil {
		xhttp.Fail(c, xerrors.New(xerrors.CodeAuthSessionExpired, "invalid credentials"))
		return
	}
	c.SetCookie(bundle.CookieName, bundle.SessionID, 7200, "/", "", protocol.CookieSameSiteLaxMode, false, true)
	xhttp.OK(c, bundle)
}

func (h *Handler) Logout(_ context.Context, c *app.RequestContext) {
	_ = xauth.DestroySession(xauth.SessionIDFromRequest(c))
	xauth.ClearSessionCookie(c)
	xhttp.OK(c, map[string]bool{"success": true})
}

func (h *Handler) Me(ctx context.Context, c *app.RequestContext) {
	user, err := h.service.GetCurrentUser(ctx, xauth.SessionIDFromRequest(c))
	if err != nil {
		xhttp.Fail(c, xerrors.New(xerrors.CodeAuthSessionExpired, "login expired"))
		return
	}
	xhttp.OK(c, user)
}
