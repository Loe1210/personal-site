package proxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
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

	handler := NewReverseProxyWithOptions(Options{
		TargetBaseURL: upstream.URL,
		StripPrefix:   "/api/content",
		MaxAttempts:   2,
	})
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

	handler := NewReverseProxyWithOptions(Options{
		TargetBaseURL: upstream.URL,
		StripPrefix:   "/api/content",
		MaxAttempts:   2,
	})
	var c app.RequestContext
	c.Request.SetRequestURI("/api/content/articles")
	c.Request.Header.SetMethod(consts.MethodPost)

	handler(context.Background(), &c)

	if attempts != 1 {
		t.Fatalf("expected post to run once, got %d attempts", attempts)
	}
}
