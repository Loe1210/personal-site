package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/pkg/xotel"
	"github.com/Loe1210/personal-site/services/content-service/internal/application"
	httpHandler "github.com/Loe1210/personal-site/services/content-service/internal/handler/http"
	infra "github.com/Loe1210/personal-site/services/content-service/internal/infra/mysql"
	articlerepo "github.com/Loe1210/personal-site/services/content-service/internal/repository/mysql"
)

var configPath = flag.String("config", "services/content-service/configs/config.yaml", "content service config path")

func main() {
	flag.Parse()
	ctx := context.Background()
	cfg, err := configs.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	shutdown, err := xotel.SetupTracerProvider(ctx, "content-service", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(ctx)
	database, err := infra.Open(cfg.MySQL)
	if err != nil {
		log.Fatal(err)
	}
	articles := application.NewArticleService(articlerepo.NewArticleRepository(database))
	handler := httpHandler.NewHandler(articles)
	h := server.Default(server.WithHostPorts(configs.GetServerAddr()))
	handler.RegisterRoutes(h)
	log.Printf("content-service listening on %s", configs.GetServerAddr())
	h.Spin()
}
