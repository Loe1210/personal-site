package site

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const githubURL = "https://github.com/Loe1210/personal-site"

type pageData struct {
	PageTitle   string
	Description string
	BodyClass   string
	Active      string
	NavClass    string
	GitHubURL   string
	Styles      []string
	Scripts     []string
	Slug        string
}

func renderPage(c *app.RequestContext, name string, data pageData) {
	c.HTML(consts.StatusOK, name, data)
}

func Home(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/home/index.html", pageData{
		PageTitle:   "Loe | Personal Site",
		Description: "A personal site about Go, Hertz, Kitex and engineering practice.",
		BodyClass:   "page-home",
		Active:      "",
		GitHubURL:   githubURL,
		Styles: []string{
			"/static/css/blog.css",
			"/static/css/home.css",
		},
		Scripts: []string{
			"/static/js/home.js",
		},
	})
}

func Blog(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/blog/index.html", pageData{
		PageTitle:   "Blog | Loe",
		Description: "Posts about Go learning, backend practice and project evolution.",
		BodyClass:   "page-blog",
		Active:      "blog",
		GitHubURL:   githubURL,
		Styles: []string{
			"/static/css/blog.css",
		},
		Scripts: []string{
			"/static/js/blog.js",
		},
	})
}

func ArticleDetail(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/article/detail.html", pageData{
		PageTitle:   "Article | Loe",
		Description: "Article detail page for the personal site.",
		BodyClass:   "page-article",
		Active:      "blog",
		GitHubURL:   githubURL,
		Styles: []string{
			"/static/css/article.css",
		},
		Scripts: []string{
			"/static/js/article.js",
		},
		Slug: c.Param("slug"),
	})
}

func About(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/about/index.html", pageData{
		PageTitle:   "About | Loe",
		Description: "About the site and the current learning path.",
		BodyClass:   "page-about",
		Active:      "about",
		GitHubURL:   githubURL,
		Styles: []string{
			"/static/css/about.css",
		},
	})
}
