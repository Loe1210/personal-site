package proxy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/Loe1210/personal-site/internal/xerrors"
	"github.com/Loe1210/personal-site/internal/xresilience"
)

func TestRewritePathStripsGatewayPrefix(t *testing.T) {
	got := RewritePath("/api/content/articles/12", "/api/content")
	if got != "/articles/12" {
		t.Fatalf("expected /articles/12, got %s", got)
	}
}

func TestReverseProxyRetriesGetOnce(t *testing.T) {
	attempts := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			http.Error(w, "temporary", http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer upstream.Close()

	handler := NewReverseProxyWithOptions(Options{TargetBaseURL: upstream.URL, StripPrefix: "/api/content", MaxAttempts: 2})
	var c app.RequestContext
	c.Request.SetRequestURI("/api/content/articles")
	c.Request.Header.SetMethod(consts.MethodGet)

	handler(context.Background(), &c)

	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
	if c.Response.StatusCode() != consts.StatusOK || string(c.Response.Body()) != "ok" {
		t.Fatalf("unexpected response: status=%d body=%s", c.Response.StatusCode(), c.Response.Body())
	}
}

func TestReverseProxyDoesNotRetryPost(t *testing.T) {
	attempts := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		http.Error(w, "temporary", http.StatusBadGateway)
	}))
	defer upstream.Close()

	handler := NewReverseProxyWithOptions(Options{TargetBaseURL: upstream.URL, StripPrefix: "/api/content", MaxAttempts: 2})
	var c app.RequestContext
	c.Request.SetRequestURI("/api/content/articles")
	c.Request.Header.SetMethod(consts.MethodPost)

	handler(context.Background(), &c)

	if attempts != 1 {
		t.Fatalf("expected post to run once, got %d attempts", attempts)
	}
	var body map[string]any
	if err := json.Unmarshal(c.Response.Body(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if int32(body["code"].(float64)) != xerrors.CodeGatewayUpstreamFailed {
		t.Fatalf("expected upstream failed envelope, got %#v", body)
	}
}

func TestReverseProxyOpenBreakerReturnsUnifiedEnvelope(t *testing.T) {
	breaker := xresilience.NewCircuitBreaker(xresilience.CircuitBreakerConfig{Name: "content", FailureThreshold: 1, OpenTimeout: time.Minute})
	breaker.RecordFailure()
	attempts := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	handler := NewReverseProxyWithOptions(Options{TargetBaseURL: upstream.URL, StripPrefix: "/api/content", Breaker: breaker})
	var c app.RequestContext
	c.Request.SetRequestURI("/api/content/articles")
	c.Request.Header.SetMethod(consts.MethodGet)

	handler(context.Background(), &c)

	if attempts != 0 {
		t.Fatalf("expected open breaker to reject before upstream call, got %d attempts", attempts)
	}
	var body map[string]any
	if err := json.Unmarshal(c.Response.Body(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if int32(body["code"].(float64)) != xerrors.CodeGatewayCircuitOpen {
		t.Fatalf("expected circuit open code, got %#v", body)
	}
	if body["data"] != nil {
		t.Fatalf("circuit open must not return fake success data: %#v", body)
	}
}
