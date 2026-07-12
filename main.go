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
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Loe1210/personal-site/biz"
	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/dal/db"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sessions"
	"github.com/hertz-contrib/sessions/cookie"
)

var configPath = flag.String("config", "configs/config.yaml", "path to config file")

func mustAbs(base string, parts ...string) string {
	all := append([]string{base}, parts...)
	return filepath.Join(all...)
}

func staticCacheMiddleware() app.HandlerFunc {
	return func(_ context.Context, c *app.RequestContext) {
		c.Next(context.Background())
		c.Response.Header.Set("Cache-Control", "public, max-age=86400")
	}
}

func serveSPA(staticRoot string) app.HandlerFunc {
	return func(_ context.Context, c *app.RequestContext) {
		reqPath := strings.TrimPrefix(c.Param("filepath"), "/")
		if reqPath == "" {
			c.File(filepath.Join(staticRoot, "index.html"))
			return
		}

		if strings.HasPrefix(reqPath, "api/") {
			c.String(consts.StatusNotFound, "Not Found")
			return
		}

		cleanPath := filepath.Clean(filepath.FromSlash(reqPath))
		target := filepath.Join(staticRoot, cleanPath)

		rel, err := filepath.Rel(staticRoot, target)
		if err != nil || strings.HasPrefix(rel, "..") {
			c.String(consts.StatusNotFound, "Not Found")
			return
		}

		if info, statErr := os.Stat(target); statErr == nil && !info.IsDir() {
			c.File(target)
			return
		}

		// SPA fallback: determine which app to serve based on path prefix
		fallback := "index.html"
		if strings.HasPrefix(reqPath, "blog/") || reqPath == "blog" {
			fallback = filepath.Join("blog", "index.html")
		} else if strings.HasPrefix(reqPath, "admin/") || reqPath == "admin" {
			fallback = filepath.Join("admin", "index.html")
		}
		c.File(filepath.Join(staticRoot, fallback))
	}
}

func main() {
	flag.Parse()

	cfg, err := configs.Load(*configPath)
	if err != nil {
		log.Printf("⚠️ Config load failed (using defaults): %v\n", err)
	}

	if err := db.Init(); err != nil {
		log.Printf("⚠️ Database init failed (API will not work): %v\n", err)
	}

	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	staticRoot := mustAbs(root, "static")

	h := server.Default(server.WithHostPorts(configs.GetServerAddr()))

	store := cookie.NewStore([]byte(cfg.Session.Secret))
	h.Use(sessions.New("personal_site_session", store))

	h.StaticFile("/swagger.json", mustAbs(root, "docs", "swagger.json"))
	h.StaticFile("/swagger.yaml", mustAbs(root, "docs", "swagger.yaml"))

	h.Use(staticCacheMiddleware())

	h.GET("/*filepath", serveSPA(staticRoot))

	biz.Register(h)

	log.Printf("🚀 Server starting on %s\n", configs.GetServerAddr())
	h.Spin()
}
