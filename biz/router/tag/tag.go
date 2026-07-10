package tag

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	taghandler "github.com/Loe1210/personal-site/biz/handler/tag"
	authmw "github.com/Loe1210/personal-site/biz/mw/session"
)

func Register(h *server.Hertz) {
	h.GET("/api/tags", taghandler.ListTags)

	admin := h.Group("/api/admin", authmw.AuthMiddleware())
	{
		admin.GET("/tags", authmw.RequirePermission("tag:read"), taghandler.ListAdminTags)
		admin.POST("/tags", authmw.RequirePermission("tag:create"), taghandler.CreateTag)
	}
}