package biz

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/services/content-service/biz/article"
	"github.com/Loe1210/personal-site/services/content-service/biz/meta"
	"github.com/Loe1210/personal-site/services/content-service/internal/service"
)

func RegisterRoutes(hertz *server.Hertz, articles *service.ArticleService, categories *service.CategoryService, tags *service.TagService) {
	article.RegisterRoutes(hertz, articles)
	meta.RegisterRoutes(hertz, categories, tags)
}
