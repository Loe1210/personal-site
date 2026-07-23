package authenticator

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"

	"github.com/Loe1210/personal-site/internal/xerrors"
	"github.com/Loe1210/personal-site/services/auth-service/internal/model"
	"github.com/Loe1210/personal-site/services/auth-service/internal/service"
)

type fakeUserRepository struct{}

func (f *fakeUserRepository) Login(context.Context, string, string) (*model.User, []string, error) {
	return nil, nil, xerrors.New(xerrors.CodeAuthSessionExpired, "invalid credentials")
}

func (f *fakeUserRepository) GetByID(context.Context, int64) (*model.User, error) {
	return nil, xerrors.New(xerrors.CodeAuthSessionExpired, "login expired")
}

func (f *fakeUserRepository) HasPermission(context.Context, int64, string) (bool, error) {
	return false, nil
}

func TestMeWithoutSessionReturnsUnifiedEnvelope(t *testing.T) {
	h := server.Default()
	handler := NewHandler(service.NewAuthService(&fakeUserRepository{}))
	h.GET("/me", func(ctx context.Context, c *app.RequestContext) {
		handler.Me(ctx, c)
	})

	resp := ut.PerformRequest(h.Engine, "GET", "/me", nil)
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
