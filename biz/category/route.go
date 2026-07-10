package category

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	sessionmw "github.com/Loe1210/personal-site/pkg/middleware/session"
)

func Register(h *server.Hertz) {
	h.GET("/api/categories", ListCategories)

	admin := h.Group("/api/admin", sessionmw.AuthMiddleware())
	{
		admin.GET("/categories", sessionmw.RequirePermission("category:read"), ListAdminCategories)
		admin.POST("/categories", sessionmw.RequirePermission("category:create"), CreateCategory)
	}
}
