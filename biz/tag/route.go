package tag

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	sessionmw "github.com/Loe1210/personal-site/pkg/middleware/session"
)

func Register(h *server.Hertz) {
	h.GET("/api/tags", ListTags)

	admin := h.Group("/api/admin", sessionmw.AuthMiddleware())
	{
		admin.GET("/tags", sessionmw.RequirePermission("tag:read"), ListAdminTags)
		admin.POST("/tags", sessionmw.RequirePermission("tag:create"), CreateTag)
	}
}
