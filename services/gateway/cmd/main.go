package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/app/server"
	kitexclient "github.com/cloudwego/kitex/client"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/kitex_gen/content/contentservice"
	contentclient "github.com/Loe1210/personal-site/services/gateway/internal/client/content"
	contenthandler "github.com/Loe1210/personal-site/services/gateway/internal/handler/content"
	"github.com/Loe1210/personal-site/services/gateway/internal/router"
	"github.com/Loe1210/personal-site/services/gateway/pkg/xnacos"
	"github.com/Loe1210/personal-site/services/gateway/pkg/xotel"
)

var configPath = flag.String("config", "configs/config.yaml", "gateway config path")

type contentRPCConfig struct {
	ServiceName string
	Address     string
	NacosAddr   string
}

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

	contentServiceClient, err := newContentServiceClient(contentRPCConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	articleClient := contentclient.NewKitexArticleClient(contentServiceClient)

	h := server.Default(server.WithHostPorts(configs.GetServerAddr()))
	deps := router.Dependencies{
		AuthServiceName: "auth-service",
		AuthBaseURL:     envOrDefault("AUTH_SERVICE_URL", "http://127.0.0.1:9001"),
		MediaBaseURL:    envOrDefault("MEDIA_SERVICE_URL", "http://127.0.0.1:9002"),
		ContentBaseURL:  envOrDefault("CONTENT_SERVICE_URL", "http://127.0.0.1:9003"),
		ContentHandler:  contenthandler.NewHandler(articleClient),
	}
	if err := router.RegisterRoutes(h, deps); err != nil {
		log.Fatal(err)
	}
	log.Printf("gateway listening on %s", configs.GetServerAddr())
	h.Spin()
}

func newContentServiceClient(cfg contentRPCConfig) (contentservice.Client, error) {
	resolver, err := xnacos.NewResolver(cfg.NacosAddr)
	if err != nil {
		return nil, err
	}
	if resolver != nil {
		return contentservice.NewClient(cfg.ServiceName, kitexclient.WithResolver(resolver))
	}
	return contentservice.NewClient(cfg.ServiceName, kitexclient.WithHostPorts(cfg.Address))
}

func contentRPCConfigFromEnv() contentRPCConfig {
	return contentRPCConfig{
		ServiceName: envOrDefault("CONTENT_SERVICE_NAME", "content-service"),
		Address:     envOrDefault("CONTENT_RPC_ADDR", "127.0.0.1:9103"),
		NacosAddr:   os.Getenv("NACOS_ADDR"),
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
