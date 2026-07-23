package article

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"

	"github.com/Loe1210/personal-site/internal/xerrors"
)

func TestGetArticleByIDInvalidIDReturnsUnifiedEnvelope(t *testing.T) {
	h := server.Default()
	handler := NewHandler(nil)
	h.GET("/articles/:id", func(ctx context.Context, c *app.RequestContext) {
		handler.GetArticleByID(ctx, c)
	})

	resp := ut.PerformRequest(h.Engine, "GET", "/articles/nope", nil)
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
}
