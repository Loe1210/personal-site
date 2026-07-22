package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/app/server"
	kitexclient "github.com/cloudwego/kitex/client"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/internal/xhttp"
	"github.com/Loe1210/personal-site/internal/xnacos"
	"github.com/Loe1210/personal-site/internal/xotel"
	"github.com/Loe1210/personal-site/internal/xsafe"
	"github.com/Loe1210/personal-site/kitex_gen/auth/authservice"
	authclient "github.com/Loe1210/personal-site/services/gateway/internal/client/auth"
	"github.com/Loe1210/personal-site/services/gateway/internal/router"
)

var configPath = flag.String("config", "configs/config.yaml", "gateway config path")

type serviceRPCConfig struct {
	ServiceName string
	Address     string
	NacosAddr   string
}

type authRPCConfig = serviceRPCConfig

func main() {
	flag.Parse()
	xsafe.InstallGoPoolPanicHandler()
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

	authServiceClient, err := newAuthServiceClient(authRPCConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	authValidator := authclient.NewKitexClient(authServiceClient)

	h := server.Default(server.WithHostPorts(configs.GetServerAddr()))
	h.Use(xhttp.Recover())
	deps := router.Dependencies{
		AuthServiceName: "auth-service",
		AuthBaseURL:     envOrDefault("AUTH_SERVICE_URL", "http://127.0.0.1:9001"),
		MediaBaseURL:    envOrDefault("MEDIA_SERVICE_URL", "http://127.0.0.1:9002"),
		ContentBaseURL:  envOrDefault("CONTENT_SERVICE_URL", "http://127.0.0.1:9003"),
		AuthValidator:   authValidator,
	}
	if err := router.RegisterRoutes(h, deps); err != nil {
		log.Fatal(err)
	}
	log.Printf("gateway listening on %s", configs.GetServerAddr())
	h.Spin()
}

func newAuthServiceClient(cfg authRPCConfig) (authservice.Client, error) {
	resolver, err := xnacos.NewResolver(cfg.NacosAddr)
	if err != nil {
		return nil, err
	}
	if resolver != nil {
		return authservice.NewClient(cfg.ServiceName, kitexclient.WithResolver(resolver))
	}
	return authservice.NewClient(cfg.ServiceName, kitexclient.WithHostPorts(cfg.Address))
}

func authRPCConfigFromEnv() authRPCConfig {
	return authRPCConfig{
		ServiceName: envOrDefault("AUTH_SERVICE_NAME", "auth-service"),
		Address:     envOrDefault("AUTH_RPC_ADDR", "127.0.0.1:9101"),
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
