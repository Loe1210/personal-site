package router

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	contenthandler "github.com/Loe1210/personal-site/services/gateway/internal/handler/content"
	"github.com/Loe1210/personal-site/services/gateway/internal/middleware"
	"github.com/Loe1210/personal-site/services/gateway/internal/proxy"
)

type Dependencies struct {
	AuthServiceName string
	AuthBaseURL     string
	ContentBaseURL  string
	MediaBaseURL    string
	BFFBaseURL      string
	ContentHandler  *contenthandler.Handler
}

func ValidateDependencies(deps Dependencies) error {
	if deps.AuthServiceName == "" {
		return errors.New("auth service name is required")
	}
	if deps.ContentHandler == nil {
		return errors.New("content handler is required")
	}
	return nil
}

func RegisterRoutes(h *server.Hertz, deps Dependencies) error {
	if err := ValidateDependencies(deps); err != nil {
		return err
	}
	h.GET("/healthz", Health)
	h.Any("/api/auth/*path", proxy.NewReverseProxy(deps.AuthBaseURL, "/api/auth"))
	uploadGuard := middleware.NewUploadGuard(middleware.UploadGuardConfig{MaxBodyBytes: 512 * 1024 * 1024, MaxConcurrent: 3, Timeout: 2 * time.Minute})
	h.Any("/api/media/*path", uploadGuard.Middleware(), proxy.NewReverseProxy(deps.MediaBaseURL, "/api/media"))
	h.GET("/api/articles", deps.ContentHandler.ListArticles)
	h.GET("/api/articles/:id", deps.ContentHandler.GetArticle)

	// Temporary compatibility path until frontend and admin APIs move to first-class gateway handlers.
	h.Any("/api/content/*path", proxy.NewReverseProxy(deps.ContentBaseURL, "/api/content"))
	if deps.BFFBaseURL != "" {
		h.Any("/api/blog/*path", proxy.NewReverseProxy(deps.BFFBaseURL, "/api/blog"))
	}
	return nil
}

func Health(_ context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, map[string]any{"status": "ok"})
}
