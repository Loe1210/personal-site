package middleware

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"

	"github.com/Loe1210/personal-site/internal/xerrors"
)

type failingSessionValidator struct{}

func (f *failingSessionValidator) ValidateSession(context.Context, string) error {
	return xerrors.New(xerrors.CodeAuthSessionExpired, "login expired")
}

func TestAuthRequiredMissingCookieReturnsUnifiedEnvelope(t *testing.T) {
	h := server.Default()
	h.GET("/admin", AuthRequired(&failingSessionValidator{}), func(ctx context.Context, c *app.RequestContext) {
		c.String(200, "ok")
	})

	resp := ut.PerformRequest(h.Engine, "GET", "/admin", nil)
	var body map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Code != 200 {
		t.Fatalf("expected http 200, got %d", resp.Code)
	}
	if body["code"].(float64) != float64(xerrors.CodeAuthLoginRequired) || body["msg"] != "login required" {
		t.Fatalf("unexpected envelope: %#v", body)
	}
}

func TestAuthRequiredValidatorErrorReturnsUnifiedEnvelope(t *testing.T) {
	h := server.Default()
	h.GET("/admin", AuthRequired(&failingSessionValidator{}), func(ctx context.Context, c *app.RequestContext) {
		c.String(200, "ok")
	})

	resp := ut.PerformRequest(h.Engine, "GET", "/admin", nil, ut.Header{Key: "Cookie", Value: "session_id=expired"})
	var body map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Code != 200 {
		t.Fatalf("expected http 200, got %d", resp.Code)
	}
	if body["code"].(float64) != float64(xerrors.CodeAuthSessionExpired) || body["msg"] != "login expired" {
		t.Fatalf("unexpected envelope: %#v", body)
	}
}
