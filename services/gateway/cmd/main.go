package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/services/gateway/internal/router"
	"github.com/Loe1210/personal-site/services/gateway/pkg/xotel"
)

var configPath = flag.String("config", "configs/config.yaml", "gateway config path")

func main() {
	flag.Parse()
	ctx := context.Background()
	_, err := configs.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	shutdown, err := xotel.SetupTracerProvider(ctx, "gateway", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(ctx)
	h := server.Default(server.WithHostPorts(configs.GetServerAddr()))
	deps := router.Dependencies{
		AuthServiceName: "auth-service",
		BFFServiceName:  "web-bff",
		AuthBaseURL:     envOrDefault("AUTH_SERVICE_URL", "http://127.0.0.1:9001"),
		MediaBaseURL:    envOrDefault("MEDIA_SERVICE_URL", "http://127.0.0.1:9002"),
		ContentBaseURL:  envOrDefault("CONTENT_SERVICE_URL", "http://127.0.0.1:9003"),
		BFFBaseURL:      envOrDefault("WEB_BFF_URL", "http://127.0.0.1:9004"),
	}
	if err := router.RegisterRoutes(h, deps); err != nil {
		log.Fatal(err)
	}
	log.Printf("gateway listening on %s", configs.GetServerAddr())
	h.Spin()
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
