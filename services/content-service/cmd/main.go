package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/services/content-service/internal/dal/db"
	kitexcontenthandler "github.com/Loe1210/personal-site/services/content-service/internal/handler/rpc"
	"github.com/Loe1210/personal-site/services/content-service/internal/service"
	"github.com/Loe1210/personal-site/services/content-service/pkg/xotel"
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
	repo := db.NewArticleRepository(database)
	articles := service.NewArticleService(repo)
	categories := service.NewCategoryService(repo)
	tags := service.NewTagService(repo)
	startContentRPCServer(cfg.RPC.Port, kitexcontenthandler.NewHandler(articles))
	h := newRouter(articles, categories, tags, configs.GetServerAddr())
	log.Printf("content-service listening on %s", configs.GetServerAddr())
	h.Spin()
}
