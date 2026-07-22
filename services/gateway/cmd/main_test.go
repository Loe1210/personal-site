package main

import "testing"

func TestAuthRPCConfigUsesDefaults(t *testing.T) {
	t.Setenv("AUTH_SERVICE_NAME", "")
	t.Setenv("AUTH_RPC_ADDR", "")
	t.Setenv("NACOS_ADDR", "")

	cfg := authRPCConfigFromEnv()

	if cfg.ServiceName != "auth-service" {
		t.Fatalf("expected default service name, got %q", cfg.ServiceName)
	}
	if cfg.Address != "127.0.0.1:9101" {
		t.Fatalf("expected default rpc address, got %q", cfg.Address)
	}
	if cfg.NacosAddr != "" {
		t.Fatalf("expected empty default nacos addr, got %q", cfg.NacosAddr)
	}
}

func TestAuthRPCConfigUsesEnvironment(t *testing.T) {
	t.Setenv("AUTH_SERVICE_NAME", "auth-service-dev")
	t.Setenv("AUTH_RPC_ADDR", "auth-service:9101")
	t.Setenv("NACOS_ADDR", "nacos:8848")

	cfg := authRPCConfigFromEnv()

	if cfg.ServiceName != "auth-service-dev" {
		t.Fatalf("expected env service name, got %q", cfg.ServiceName)
	}
	if cfg.Address != "auth-service:9101" {
		t.Fatalf("expected env rpc address, got %q", cfg.Address)
	}
	if cfg.NacosAddr != "nacos:8848" {
		t.Fatalf("expected env nacos addr, got %q", cfg.NacosAddr)
	}
}
