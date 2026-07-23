package meta

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"

	"github.com/Loe1210/personal-site/internal/xerrors"
	"github.com/Loe1210/personal-site/services/content-service/internal/model"
	"github.com/Loe1210/personal-site/services/content-service/internal/service"
)

type fakeMetaRepo struct {
	listCategoriesErr error
}

func (f *fakeMetaRepo) ListCategories(context.Context) ([]model.Category, error) {
	if f.listCategoriesErr != nil {
		return nil, f.listCategoriesErr
	}
	return []model.Category{{ID: 1, Name: "Go", Slug: "go"}}, nil
}

func (f *fakeMetaRepo) CreateCategory(context.Context, *model.Category) error { return nil }
func (f *fakeMetaRepo) UpdateCategory(context.Context, *model.Category) error { return nil }
func (f *fakeMetaRepo) DeleteCategory(context.Context, int64) error           { return nil }
func (f *fakeMetaRepo) ListTags(context.Context) ([]model.Tag, error)         { return nil, nil }
func (f *fakeMetaRepo) CreateTag(context.Context, *model.Tag) error           { return nil }
func (f *fakeMetaRepo) UpdateTag(context.Context, *model.Tag) error           { return nil }
func (f *fakeMetaRepo) DeleteTag(context.Context, int64) error                { return nil }

func TestUpdateCategoryInvalidIDReturnsUnifiedEnvelope(t *testing.T) {
	h := server.Default()
	handler := NewHandler(nil, nil)
	h.PUT("/categories/:id", func(ctx context.Context, c *app.RequestContext) {
		handler.UpdateCategory(ctx, c)
	})

	resp := ut.PerformRequest(h.Engine, "PUT", "/categories/nope", nil)
	body := decodeEnvelope(t, resp.Body.Bytes())

	if resp.Code != 200 {
		t.Fatalf("expected http 200, got %d", resp.Code)
	}
	if body["code"].(float64) != float64(xerrors.CodeInvalidArgument) || body["msg"] != "invalid category id" {
		t.Fatalf("unexpected envelope: %#v", body)
	}
	if _, ok := body["message"]; ok {
		t.Fatalf("expected msg field only, got legacy message in %#v", body)
	}
}

func TestListCategoriesServiceErrorReturnsUnifiedEnvelope(t *testing.T) {
	repo := &fakeMetaRepo{listCategoriesErr: errors.New("database failed")}
	handler := NewHandler(service.NewCategoryService(repo), service.NewTagService(repo))
	h := server.Default()
	h.GET("/categories", func(ctx context.Context, c *app.RequestContext) {
		handler.ListCategories(ctx, c)
	})

	resp := ut.PerformRequest(h.Engine, "GET", "/categories", nil)
	body := decodeEnvelope(t, resp.Body.Bytes())

	if resp.Code != 200 {
		t.Fatalf("expected http 200, got %d", resp.Code)
	}
	if body["code"].(float64) != float64(xerrors.CodeInternal) || body["msg"] != "internal error" {
		t.Fatalf("unexpected envelope: %#v", body)
	}
	if _, ok := body["message"]; ok {
		t.Fatalf("expected msg field only, got legacy message in %#v", body)
	}
}

func decodeEnvelope(t *testing.T, payload []byte) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(payload, &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	return body
}
