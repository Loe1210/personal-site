package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/pkg/xotel"
	db "github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
	"github.com/Loe1210/personal-site/services/media-service/internal/dal/storage"
	"github.com/Loe1210/personal-site/services/media-service/internal/service"
)

var configPath = flag.String("config", "services/media-service/configs/config.yaml", "media service config path")

func main() {
	flag.Parse()
	ctx := context.Background()
	cfg, err := configs.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	shutdown, err := xotel.SetupTracerProvider(ctx, "media-service", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
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
	store := storage.NewLocalStorage(cfg.Upload.RootDir, cfg.Upload.PublicBasePath)
	media := service.NewMediaService(store, db.NewFileRepository(database))
	h := newRouter(media, configs.GetServerAddr())
	log.Printf("media-service listening on %s", configs.GetServerAddr())
	h.Spin()
}