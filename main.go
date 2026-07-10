// @title Personal Site API
// @version 1.0
// @description Personal site backend API for learning Hertz, Kitex and future evolution.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sessions"
	"github.com/hertz-contrib/sessions/cookie"

	"github.com/Loe1210/personal-site/biz/dal/db"
	"github.com/Loe1210/personal-site/biz/router"
)

func mustAbs(base string, parts ...string) string {
	all := append([]string{base}, parts...)
	return filepath.Join(all...)
}

func staticFileHandler(staticRoot string) app.HandlerFunc {
	return func(_ context.Context, c *app.RequestContext) {
		reqPath := strings.TrimPrefix(c.Param("filepath"), "/")
		cleanPath := filepath.Clean(filepath.FromSlash(reqPath))
		target := filepath.Join(staticRoot, cleanPath)

		rel, err := filepath.Rel(staticRoot, target)
		if err != nil || strings.HasPrefix(rel, "..") {
			c.String(consts.StatusNotFound, "Not Found")
			return
		}

		if _, err := os.Stat(target); err != nil {
			c.String(consts.StatusNotFound, "Not Found")
			return
		}

		c.File(target)
	}
}

func main() {
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	staticRoot := mustAbs(root, "static")

	h := server.Default()
	h.LoadHTMLFiles(
		mustAbs(root, "templates", "components", "layout.html"),
		mustAbs(root, "templates", "components", "nav.html"),
		mustAbs(root, "templates", "components", "article-card.html"),
		mustAbs(root, "templates", "pages", "home", "index.html"),
		mustAbs(root, "templates", "pages", "blog", "index.html"),
		mustAbs(root, "templates", "pages", "article", "detail.html"),
		mustAbs(root, "templates", "pages", "about", "index.html"),
	)

	store := cookie.NewStore([]byte("personal-site-session-secret"))
	h.Use(sessions.New("personal_site_session", store))

	h.StaticFile("/swagger.json", mustAbs(root, "docs", "swagger.json"))
	h.StaticFile("/swagger.yaml", mustAbs(root, "docs", "swagger.yaml"))
	h.GET("/static/*filepath", staticFileHandler(staticRoot))
	h.HEAD("/static/*filepath", staticFileHandler(staticRoot))

	router.Register(h)

	h.Spin()
}
