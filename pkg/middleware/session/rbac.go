package session

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	dbmodel "github.com/Loe1210/personal-site/dal/db"
	"github.com/Loe1210/personal-site/pkg/errno"
	"github.com/Loe1210/personal-site/pkg/response"
	"github.com/Loe1210/personal-site/pkg/xauth"
)

func RequirePermission(permissionCode string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		claims, exists := xauth.ClaimsFromContext(c)
		if !exists {
			response.WriteErrorMessage(c, errno.Unauthorized, "login required")
			c.Abort()
			return
		}

		allowed, err := hasPermissionForRoles(claims.Roles, permissionCode)
		if err != nil {
			response.WriteError(c, errno.Internal)
			c.Abort()
			return
		}
		if !allowed {
			response.WriteError(c, errno.Forbidden)
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

func hasPermissionForRoles(roleCodes []string, permissionCode string) (bool, error) {
	if len(roleCodes) == 0 {
		return false, nil
	}

	var count int64
	err := dbmodel.DB.
		Table("roles").
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("roles.code IN ? AND permissions.code = ?", roleCodes, permissionCode).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
