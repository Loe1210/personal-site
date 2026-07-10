package article

import (
	"github.com/cloudwego/hertz/pkg/app"

	articlemodel "github.com/Loe1210/personal-site/biz/model/article"
)

func bindListRequest(c *app.RequestContext) (*articlemodel.ListArticlesRequest, error) {
	var req articlemodel.ListArticlesRequest
	if err := c.BindAndValidate(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
