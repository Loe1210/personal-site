package http

import (
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/Loe1210/personal-site/services/content-service/internal/service"
)

func listFilterFromRequest(c *app.RequestContext) service.ListFilter {
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.Query("page_size"), 10, 64)
	return service.ListFilter{
		Page:     page,
		PageSize: pageSize,
		Status:   c.Query("status"),
		Keyword:  c.Query("keyword"),
	}
}
