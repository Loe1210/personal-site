package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/pkg/xotel"
	"github.com/Loe1210/personal-site/services/content-service/internal/dal/db"
	"github.com/Loe1210/personal-site/services/content-service/internal/service"
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
	database, err := db.Open(cfg.MySQL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}
	articles := service.NewArticleService(db.NewArticleRepository(database))
	h := newRouter(articles, configs.GetServerAddr())
	log.Printf("content-service listening on %s", configs.GetServerAddr())
	h.Spin()
}