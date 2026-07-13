package auth

import (
	"context"

	authmodel "github.com/Loe1210/personal-site/biz/model/auth"
	"github.com/Loe1210/personal-site/dal/db"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
	"github.com/Loe1210/personal-site/pkg/xauth"
	"github.com/Loe1210/personal-site/pkg/xtrace"
	authservice "github.com/Loe1210/personal-site/service"
	"github.com/cloudwego/hertz/pkg/app"
)

type loginResponse struct {
	User          *authmodel.User      `json:"user"`
	SessionBundle *xauth.SessionBundle `json:"session_bundle"`
}

func Login(ctx context.Context, c *app.RequestContext) {
	req, err := bindLoginRequest(c)
	if err != nil {
		response.WriteError(c, errno.BadRequest)
		return
	}

	resp, err := authservice.Login(ctx, req)
	if err != nil {
		if appErr, ok := err.(*errno.AppError); ok {
			response.WriteError(c, appErr)
			return
		}
		response.WriteError(c, errno.Internal)
		return
	}

	roles, err := loadUserRoles(resp.User.ID)
	if err != nil {
		response.WriteError(c, errno.Internal)
		return
	}

	traceID := xtrace.EnsureTraceID(c)
	bundle, err := xauth.CreateSessionBundleWithTrace(resp.User.ID, resp.User.Username, roles, traceID)
	if err != nil {
		response.WriteError(c, errno.Internal)
		return
	}

	xauth.WriteSessionCookie(c, bundle)
	response.WriteSuccess(c, &loginResponse{User: resp.User, SessionBundle: bundle})
}

func Me(_ context.Context, c *app.RequestContext) {
	claims, ok := xauth.ClaimsFromContext(c)
	if !ok {
		response.WriteError(c, errno.Unauthorized)
		return
	}

	var user db.User
	if err := db.DB.First(&user, claims.UserID).Error; err != nil {
		response.WriteError(c, errno.Unauthorized)
		return
	}
	if user.Username != claims.Username {
		response.WriteError(c, errno.Unauthorized)
		return
	}

	resp := &authmodel.GetCurrentUserResponse{
		User: &authmodel.User{
			ID:        user.ID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	response.WriteSuccess(c, resp)
}

func Logout(_ context.Context, c *app.RequestContext) {
	sessionID := xauth.SessionIDFromRequest(c)
	_ = xauth.DestroySession(sessionID)
	xauth.ClearSessionCookie(c)

	response.WriteSuccess(c, &authmodel.LogoutResponse{
		Message: "logout success",
	})
}

func loadUserRoles(userID int64) ([]string, error) {
	var roles []string
	if err := db.DB.
		Table("roles").
		Select("roles.code").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Order("roles.id ASC").
		Scan(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
