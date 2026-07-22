package main

import "testing"

func TestContentRPCConfigUsesDefaults(t *testing.T) {
	t.Setenv("CONTENT_SERVICE_NAME", "")
	t.Setenv("CONTENT_RPC_ADDR", "")
	t.Setenv("NACOS_ADDR", "")

	cfg := contentRPCConfigFromEnv()

	if cfg.ServiceName != "content-service" {
		t.Fatalf("expected default service name, got %q", cfg.ServiceName)
	}
	if cfg.Address != "127.0.0.1:9103" {
		t.Fatalf("expected default rpc address, got %q", cfg.Address)
	}
	if cfg.NacosAddr != "" {
		t.Fatalf("expected empty default nacos addr, got %q", cfg.NacosAddr)
	}
}

func TestContentRPCConfigUsesEnvironment(t *testing.T) {
	t.Setenv("CONTENT_SERVICE_NAME", "content-service-dev")
	t.Setenv("CONTENT_RPC_ADDR", "content-service:9103")
	t.Setenv("NACOS_ADDR", "nacos:8848")

	cfg := contentRPCConfigFromEnv()

	if cfg.ServiceName != "content-service-dev" {
		t.Fatalf("expected env service name, got %q", cfg.ServiceName)
	}
	if cfg.Address != "content-service:9103" {
		t.Fatalf("expected env rpc address, got %q", cfg.Address)
	}
	if cfg.NacosAddr != "nacos:8848" {
		t.Fatalf("expected env nacos addr, got %q", cfg.NacosAddr)
	}
}
