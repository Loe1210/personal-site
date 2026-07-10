package site

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const githubURL = "https://github.com/Loe1210/personal-site"

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
</head>
<body style="margin:0;">
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      SwaggerUIBundle({
        url: '/swagger.json',
        dom_id: '#swagger-ui'
      });
    };
  </script>
</body>
</html>`

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
	ArticleID   string
}

func renderPage(c *app.RequestContext, name string, data pageData) {
	c.HTML(consts.StatusOK, name, data)
}

func Home(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/home/index.html", pageData{
		PageTitle:   "Loe | Personal Site",
		Description: "A personal site about Go, Hertz, Kitex and engineering practice.",
		BodyClass:   "page-home",
		GitHubURL:   githubURL,
		Styles:      []string{"/static/css/blog.css", "/static/css/home.css"},
		Scripts:     []string{"/static/js/home.js"},
	})
}

func Blog(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/blog/index.html", pageData{
		PageTitle:   "Blog | Loe",
		Description: "Posts about Go learning, backend practice and project evolution.",
		BodyClass:   "page-blog",
		Active:      "blog",
		GitHubURL:   githubURL,
		Styles:      []string{"/static/css/blog.css"},
		Scripts:     []string{"/static/js/blog.js"},
	})
}

func ArticleDetail(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/article/detail.html", pageData{
		PageTitle:   "Article | Loe",
		Description: "Article detail page for the personal site.",
		BodyClass:   "page-article",
		Active:      "blog",
		GitHubURL:   githubURL,
		Styles:      []string{"/static/css/article.css"},
		Scripts:     []string{"/static/js/article.js"},
		Slug:        c.Param("slug"),
	})
}

func About(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/about/index.html", pageData{
		PageTitle:   "About | Loe",
		Description: "About the site and the current learning path.",
		BodyClass:   "page-about",
		Active:      "about",
		GitHubURL:   githubURL,
		Styles:      []string{"/static/css/about.css"},
	})
}

func AdminLogin(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/admin/login.html", pageData{
		PageTitle:   "Admin Login | Loe",
		Description: "Admin login for the personal site.",
		BodyClass:   "page-admin-login",
		GitHubURL:   githubURL,
		Styles:      []string{"/static/css/admin.css", "/static/css/admin-login.css"},
		Scripts:     []string{"/static/js/admin-common.js", "/static/js/admin-login.js"},
	})
}

func AdminDashboard(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/admin/dashboard.html", pageData{
		PageTitle:   "Dashboard | Loe Admin",
		Description: "Admin dashboard for content operations.",
		BodyClass:   "page-admin dashboard-page",
		Active:      "dashboard",
		GitHubURL:   githubURL,
		Styles:      []string{"/static/css/admin.css", "/static/css/admin-dashboard.css"},
		Scripts:     []string{"/static/js/admin-common.js", "/static/js/admin-dashboard.js"},
	})
}

func AdminArticles(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/admin/articles.html", pageData{
		PageTitle:   "Articles | Loe Admin",
		Description: "Admin article management.",
		BodyClass:   "page-admin articles-page",
		Active:      "articles",
		GitHubURL:   githubURL,
		Styles:      []string{"/static/css/admin.css", "/static/css/admin-articles.css"},
		Scripts:     []string{"/static/js/admin-common.js", "/static/js/admin-articles.js"},
	})
}

func AdminArticleNew(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/admin/article-edit.html", pageData{
		PageTitle:   "New Article | Loe Admin",
		Description: "Create a new article.",
		BodyClass:   "page-admin editor-page",
		Active:      "articles",
		GitHubURL:   githubURL,
		Styles:      []string{"/static/css/admin.css", "/static/css/admin-editor.css"},
		Scripts:     []string{"/static/js/admin-common.js", "/static/js/admin-editor.js"},
	})
}

func AdminArticleEdit(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/admin/article-edit.html", pageData{
		PageTitle:   "Edit Article | Loe Admin",
		Description: "Edit an existing article.",
		BodyClass:   "page-admin editor-page",
		Active:      "articles",
		GitHubURL:   githubURL,
		Styles:      []string{"/static/css/admin.css", "/static/css/admin-editor.css"},
		Scripts:     []string{"/static/js/admin-common.js", "/static/js/admin-editor.js"},
		ArticleID:   c.Param("id"),
	})
}

func AdminTaxonomy(_ context.Context, c *app.RequestContext) {
	renderPage(c, "pages/admin/taxonomy.html", pageData{
		PageTitle:   "Taxonomy | Loe Admin",
		Description: "Manage categories and tags.",
		BodyClass:   "page-admin taxonomy-page",
		Active:      "taxonomy",
		GitHubURL:   githubURL,
		Styles:      []string{"/static/css/admin.css", "/static/css/admin-taxonomy.css"},
		Scripts:     []string{"/static/js/admin-common.js", "/static/js/admin-taxonomy.js"},
	})
}

func Swagger(_ context.Context, c *app.RequestContext) {
	c.Data(consts.StatusOK, "text/html; charset=utf-8", []byte(swaggerHTML))
}
