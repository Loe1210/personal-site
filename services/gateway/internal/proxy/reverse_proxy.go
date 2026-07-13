package proxy

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func RewritePath(path string, stripPrefix string) string {
	rewritten := strings.TrimPrefix(path, stripPrefix)
	if rewritten == "" {
		return "/"
	}
	if !strings.HasPrefix(rewritten, "/") {
		return "/" + rewritten
	}
	return rewritten
}

func NewReverseProxy(targetBaseURL string, stripPrefix string) app.HandlerFunc {
	client := &http.Client{Timeout: 10 * time.Second}
	baseURL := strings.TrimRight(targetBaseURL, "/")
	return func(ctx context.Context, c *app.RequestContext) {
		if baseURL == "" {
			c.JSON(consts.StatusBadGateway, map[string]any{"code": 40001, "message": "upstream is not configured"})
			return
		}
		path := RewritePath(string(c.Path()), stripPrefix)
		if query := string(c.QueryArgs().QueryString()); query != "" {
			path += "?" + query
		}
		req, err := http.NewRequestWithContext(ctx, string(c.Method()), baseURL+path, strings.NewReader(string(c.Request.Body())))
		if err != nil {
			c.JSON(consts.StatusBadGateway, map[string]any{"code": 40002, "message": "build upstream request failed"})
			return
		}
		c.Request.Header.VisitAll(func(key, value []byte) {
			req.Header.Set(string(key), string(value))
		})
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(consts.StatusBadGateway, map[string]any{"code": 40003, "message": "upstream request failed"})
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(consts.StatusBadGateway, map[string]any{"code": 40004, "message": "read upstream response failed"})
			return
		}
		for key, values := range resp.Header {
			for _, value := range values {
				c.Response.Header.Add(key, value)
			}
		}
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}
}
