package router

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/Loe1210/personal-site/services/gateway/internal/middleware"
	"github.com/Loe1210/personal-site/services/gateway/internal/proxy"
)

type Dependencies struct {
	AuthServiceName string
	BFFServiceName  string
	AuthBaseURL     string
	ContentBaseURL  string
	MediaBaseURL    string
	BFFBaseURL      string
}

func ValidateDependencies(deps Dependencies) error {
	if deps.AuthServiceName == "" {
		return errors.New("auth service name is required")
	}
	if deps.BFFServiceName == "" {
		return errors.New("bff service name is required")
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
	h.Any("/api/content/*path", proxy.NewReverseProxy(deps.ContentBaseURL, "/api/content"))
	h.Any("/api/blog/*path", proxy.NewReverseProxy(deps.BFFBaseURL, "/api/blog"))
	return nil
}

func Health(_ context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, map[string]any{"status": "ok"})
}
