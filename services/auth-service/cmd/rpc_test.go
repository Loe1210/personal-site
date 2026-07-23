package main

import "testing"

func TestAuthServiceRPCConfigUsesDefaults(t *testing.T) {
	t.Setenv("SERVICE_NAME", "")
	t.Setenv("NACOS_ADDR", "")

	cfg := authServiceRPCConfigFromEnv()

	if cfg.ServiceName != "auth-service" {
		t.Fatalf("expected default service name, got %q", cfg.ServiceName)
	}
	if cfg.NacosAddr != "" {
		t.Fatalf("expected empty nacos addr, got %q", cfg.NacosAddr)
	}
}

func TestAuthServiceRPCConfigUsesEnvironment(t *testing.T) {
	t.Setenv("SERVICE_NAME", "auth-service-dev")
	t.Setenv("NACOS_ADDR", "nacos:8848")

	cfg := authServiceRPCConfigFromEnv()

	if cfg.ServiceName != "auth-service-dev" {
		t.Fatalf("expected env service name, got %q", cfg.ServiceName)
	}
	if cfg.NacosAddr != "nacos:8848" {
		t.Fatalf("expected env nacos addr, got %q", cfg.NacosAddr)
	}
}
