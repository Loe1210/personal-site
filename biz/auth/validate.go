package auth

import (
	"github.com/cloudwego/hertz/pkg/app"

	authmodel "github.com/Loe1210/personal-site/biz/model/auth"
)

func bindLoginRequest(c *app.RequestContext) (*authmodel.UserLoginRequest, error) {
	var req authmodel.UserLoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
