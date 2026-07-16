package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

var (
	ErrUploadTooLarge = errors.New("upload request too large")
	ErrUploadBusy     = errors.New("too many concurrent uploads")
)

type UploadGuardConfig struct {
	MaxBodyBytes  int64
	MaxConcurrent int
	Timeout       time.Duration
}

type UploadGuard struct {
	maxBodyBytes int64
	sem          chan struct{}
	timeout      time.Duration
}

func NewUploadGuard(cfg UploadGuardConfig) *UploadGuard {
	maxConcurrent := cfg.MaxConcurrent
	if maxConcurrent <= 0 {
		maxConcurrent = 3
	}
	return &UploadGuard{maxBodyBytes: cfg.MaxBodyBytes, sem: make(chan struct{}, maxConcurrent), timeout: cfg.Timeout}
}

func (g *UploadGuard) Acquire(contentLength int64) (func(), error) {
	if g == nil {
		return func() {}, nil
	}
	if g.maxBodyBytes > 0 && contentLength > g.maxBodyBytes {
		return nil, ErrUploadTooLarge
	}
	select {
	case g.sem <- struct{}{}:
		return func() { <-g.sem }, nil
	default:
		return nil, ErrUploadBusy
	}
}

func (g *UploadGuard) Middleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		contentLength := int64(c.Request.Header.ContentLength())
		release, err := g.Acquire(contentLength)
		if err != nil {
			status := consts.StatusTooManyRequests
			code := 40012
			if errors.Is(err, ErrUploadTooLarge) {
				status = consts.StatusRequestEntityTooLarge
				code = 40011
			}
			c.JSON(status, map[string]any{"code": code, "message": err.Error()})
			c.Abort()
			return
		}
		defer release()
		if g.timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, g.timeout)
			defer cancel()
		}
		c.Next(ctx)
	}
}
