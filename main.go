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
		path := string(c.Request.URI().Path())
		if strings.HasPrefix(path, "/api/") || path == "/health" {
			c.Response.Header.Set("Cache-Control", "no-store")
		} else {
			c.Response.Header.Set("Cache-Control", "public, max-age=86400")
		}
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

		// If a path like blog/post/:id is requested, try to serve blog/post.html
		if htmlFallback := tryHTMLPage(staticRoot, reqPath); htmlFallback != "" {
			c.File(htmlFallback)
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

func tryHTMLPage(staticRoot, reqPath string) string {
	for p := reqPath; p != "." && p != "/" && p != ""; p = filepath.Dir(p) {
		if p == "." {
			break
		}
		candidate := filepath.Join(staticRoot, p+".html")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate
		}
	}
	return ""
}

func main() {
	flag.Parse()

	cfg, err := configs.Load(*configPath)
	if err != nil {
		log.Printf("Config load failed (using defaults): %v\n", err)
	}

	if err := db.Init(); err != nil {
		log.Printf("Database init failed (API will not work): %v\n", err)
	}

	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	staticRoot := mustAbs(root, "static")

	h := server.Default(server.WithHostPorts(configs.GetServerAddr()))

	store := cookie.NewStore([]byte(cfg.Session.Secret))
	h.Use(sessions.New("personal_site_session", store))

	h.Use(staticCacheMiddleware())

	biz.Register(h)

	h.GET("/*filepath", serveSPA(staticRoot))

	log.Printf("Server starting on %s\n", configs.GetServerAddr())
	h.Spin()
}
