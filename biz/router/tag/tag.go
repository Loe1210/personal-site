package tag

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	taghandler "github.com/Loe1210/personal-site/biz/handler/tag"
	"github.com/Loe1210/personal-site/biz/mw/jwt"
)

func Register(h *server.Hertz) {
	h.GET("/api/tags", taghandler.ListTags)

	admin := h.Group("/api/admin", jwt.AuthMiddleware())
	{
		admin.GET("/tags", taghandler.ListAdminTags)
		admin.POST("/tags", taghandler.CreateTag)
	}
}