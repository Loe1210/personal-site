package main

import "testing"

func TestContentServiceRPCConfigUsesDefaults(t *testing.T) {
	t.Setenv("SERVICE_NAME", "")
	t.Setenv("NACOS_ADDR", "")

	cfg := contentServiceRPCConfigFromEnv()

	if cfg.ServiceName != "content-service" {
		t.Fatalf("expected default service name, got %q", cfg.ServiceName)
	}
	if cfg.NacosAddr != "" {
		t.Fatalf("expected empty nacos addr, got %q", cfg.NacosAddr)
	}
}

func TestContentServiceRPCConfigUsesEnvironment(t *testing.T) {
	t.Setenv("SERVICE_NAME", "content-service-dev")
	t.Setenv("NACOS_ADDR", "nacos:8848")

	cfg := contentServiceRPCConfigFromEnv()

	if cfg.ServiceName != "content-service-dev" {
		t.Fatalf("expected env service name, got %q", cfg.ServiceName)
	}
	if cfg.NacosAddr != "nacos:8848" {
		t.Fatalf("expected env nacos addr, got %q", cfg.NacosAddr)
	}
}
