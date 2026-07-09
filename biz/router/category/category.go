package category

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	categoryhandler "github.com/Loe1210/personal-site/biz/handler/category"
	"github.com/Loe1210/personal-site/biz/mw/jwt"
)

func Register(h *server.Hertz) {
	h.GET("/api/categories", categoryhandler.ListCategories)

	admin := h.Group("/api/admin", jwt.AuthMiddleware())
	{
		admin.GET("/categories", categoryhandler.ListAdminCategories)
		admin.POST("/categories", categoryhandler.CreateCategory)
	}
}