package auth

import (
	"context"
	"errors"
	"time"

	authjwt "github.com/Loe1210/personal-site/biz/mw/jwt"
	authmodel "github.com/Loe1210/personal-site/biz/model/auth"
)

const (
	adminUsername = "admin"
	adminPassword = "123456"
	adminUserID   = int64(1)
)

func Login(_ context.Context, req *authmodel.AdminLoginRequest) (*authmodel.AdminLoginResponse, error) {
	if req.Username != adminUsername || req.Password != adminPassword {
		return nil, errors.New("invalid username or password")
	}

	token, err := authjwt.GenerateToken(adminUserID, adminUsername)
	if err != nil {
		return nil, err
	}

	return &authmodel.AdminLoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		User: &authmodel.AdminUser{
			Id:       adminUserID,
			Username: adminUsername,
			Nickname: "Administrator",
		},
	}, nil
}