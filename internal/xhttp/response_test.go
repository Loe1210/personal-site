package xhttp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"

	"github.com/Loe1210/personal-site/internal/xerrors"
)

func TestOKWritesUnifiedEnvelope(t *testing.T) {
	h := server.Default()
	h.GET("/ok", func(ctx context.Context, c *app.RequestContext) {
		OK(c, map[string]string{"name": "article"})
	})

	resp := ut.PerformRequest(h.Engine, "GET", "/ok", nil)
	var body map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Code != 200 {
		t.Fatalf("expected http 200, got %d", resp.Code)
	}
	if body["code"].(float64) != 0 || body["msg"] != "success" {
		t.Fatalf("unexpected envelope: %#v", body)
	}
	if _, ok := body["data"].(map[string]any); !ok {
		t.Fatalf("expected data object, got %#v", body["data"])
	}
}

func TestFailWritesUnifiedEnvelopeWithHTTP200(t *testing.T) {
	h := server.Default()
	h.GET("/fail", func(ctx context.Context, c *app.RequestContext) {
		Fail(c, xerrors.New(xerrors.CodeInvalidArgument, "invalid article id"))
	})

	resp := ut.PerformRequest(h.Engine, "GET", "/fail", nil)
	var body map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Code != 200 {
		t.Fatalf("expected http 200, got %d", resp.Code)
	}
	if body["code"].(float64) != float64(xerrors.CodeInvalidArgument) || body["msg"] != "invalid article id" {
		t.Fatalf("unexpected envelope: %#v", body)
	}
	if body["data"] != nil {
		t.Fatalf("expected nil data for error, got %#v", body["data"])
	}
}

func TestRecoverConvertsPanicToUnifiedEnvelope(t *testing.T) {
	h := server.New()
	h.Use(Recover())
	h.GET("/panic", func(ctx context.Context, c *app.RequestContext) {
		panic("boom")
	})

	resp := ut.PerformRequest(h.Engine, "GET", "/panic", nil)
	var body map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Code != 200 {
		t.Fatalf("expected http 200, got %d", resp.Code)
	}
	if body["code"].(float64) != float64(xerrors.CodeInternal) || body["msg"] != "internal error" {
		t.Fatalf("unexpected recover envelope: %#v", body)
	}
}
