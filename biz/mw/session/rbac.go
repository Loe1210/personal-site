package session

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	dbmodel "github.com/Loe1210/personal-site/biz/dal/db"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
)

func RequirePermission(permissionCode string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userIDValue, exists := c.Get("user_id")
		if !exists {
			response.WriteErrorMessage(c, errno.Unauthorized, "login required")
			c.Abort()
			return
		}

		userID, ok := userIDValue.(int64)
		if !ok {
			response.WriteErrorMessage(c, errno.Unauthorized, "invalid user session")
			c.Abort()
			return
		}

		var count int64
		err := dbmodel.DB.
			Table("user_roles").
			Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").
			Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
			Where("user_roles.user_id = ? AND permissions.code = ?", userID, permissionCode).
			Count(&count).Error
		if err != nil {
			response.WriteError(c, errno.Internal)
			c.Abort()
			return
		}

		if count == 0 {
			response.WriteError(c, errno.Forbidden)
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}