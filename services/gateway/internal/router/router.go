package router

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Dependencies struct {
	AuthServiceName string
	BFFServiceName  string
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
	return nil
}

func Health(_ context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, map[string]any{"status": "ok"})
}
