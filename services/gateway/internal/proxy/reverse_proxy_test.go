package proxy

import "testing"

func TestRewritePathStripsGatewayPrefix(t *testing.T) {
	got := RewritePath("/api/content/articles/12", "/api/content")
	if got != "/articles/12" {
		t.Fatalf("expected /articles/12, got %s", got)
	}
}
