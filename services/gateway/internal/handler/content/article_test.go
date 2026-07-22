package content

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	contentclient "github.com/Loe1210/personal-site/services/gateway/internal/client/content"
)

type fakeArticleClient struct {
	listReq contentclient.ListPublicArticlesRequest
	listErr error
}

func (f *fakeArticleClient) ListPublicArticles(_ context.Context, req contentclient.ListPublicArticlesRequest) (*contentclient.ListPublicArticlesResponse, error) {
	f.listReq = req
	if f.listErr != nil {
		return nil, f.listErr
	}
	return &contentclient.ListPublicArticlesResponse{
		List: []contentclient.Article{
			{ID: 7, Title: "Gateway RPC", Slug: "gateway-rpc", Summary: "summary"},
		},
		Total: 1,
	}, nil
}

func (f *fakeArticleClient) GetArticleByID(_ context.Context, id int64) (*contentclient.Article, error) {
	return &contentclient.Article{ID: id, Title: "Gateway RPC", Slug: "gateway-rpc"}, nil
}

func TestListArticlesReturnsClientData(t *testing.T) {
	client := &fakeArticleClient{}
	handler := NewHandler(client)
	ctx := ut.CreateUtRequestContext("GET", "/api/articles?page=2&page_size=20&keyword=go", nil)

	handler.ListArticles(context.Background(), ctx)

	if ctx.Response.StatusCode() != consts.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", ctx.Response.StatusCode(), string(ctx.Response.Body()))
	}
	if client.listReq.Page != 2 || client.listReq.PageSize != 20 || client.listReq.Keyword != "go" {
		t.Fatalf("unexpected list request: %#v", client.listReq)
	}
	var body contentclient.ListPublicArticlesResponse
	if err := json.Unmarshal(ctx.Response.Body(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Total != 1 || len(body.List) != 1 || body.List[0].Title != "Gateway RPC" {
		t.Fatalf("unexpected response body: %#v", body)
	}
}

func TestListArticlesReturnsBadGatewayOnClientError(t *testing.T) {
	handler := NewHandler(&fakeArticleClient{listErr: errors.New("content unavailable")})
	ctx := ut.CreateUtRequestContext("GET", "/api/articles", nil)

	handler.ListArticles(context.Background(), ctx)

	if ctx.Response.StatusCode() != consts.StatusBadGateway {
		t.Fatalf("expected status 502, got %d", ctx.Response.StatusCode())
	}
}

func TestGetArticleRejectsInvalidID(t *testing.T) {
	handler := NewHandler(&fakeArticleClient{})
	ctx := ut.CreateUtRequestContext("GET", "/api/articles/not-a-number", nil)

	handler.GetArticle(context.Background(), ctx)

	if ctx.Response.StatusCode() != consts.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", ctx.Response.StatusCode())
	}
}
