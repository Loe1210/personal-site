package category

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	categoryhandler "github.com/Loe1210/personal-site/biz/handler/category"
	authmw "github.com/Loe1210/personal-site/biz/mw/session"
)

func Register(h *server.Hertz) {
	h.GET("/api/categories", categoryhandler.ListCategories)

	admin := h.Group("/api/admin", authmw.AuthMiddleware())
	{
		admin.GET("/categories", authmw.RequirePermission("category:read"), categoryhandler.ListAdminCategories)
		admin.POST("/categories", authmw.RequirePermission("category:create"), categoryhandler.CreateCategory)
	}
}