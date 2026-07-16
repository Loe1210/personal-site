package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/Loe1210/personal-site/configs"
	db "github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
	"github.com/Loe1210/personal-site/services/media-service/internal/dal/storage"
	kitexmediahandler "github.com/Loe1210/personal-site/services/media-service/internal/handler/rpc"
	"github.com/Loe1210/personal-site/services/media-service/internal/service"
	"github.com/Loe1210/personal-site/services/media-service/pkg/xotel"
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
	fileRepo := db.NewFileRepository(database)
	uploadTaskRepo := db.NewUploadTaskRepository(database)
	uploadChunkRepo := db.NewUploadChunkRepository(database)
	uploadTasks := service.NewUploadTaskService(&cfg.Upload, uploadTaskRepo, uploadChunkRepo)
	uploadTasks.ConfigureCompletion(service.NewMergeService(cfg.Upload.TmpRootDir, cfg.Upload.RootDir, cfg.Upload.PublicBasePath), fileRepo, service.NewImageProcessor())
	tmpStore := storage.NewTmpStorage(cfg.Upload.TmpRootDir)
	chunks := service.NewChunkService(uploadTaskRepo, uploadChunkRepo, tmpStore)
	media := service.NewMediaService(store, fileRepo)
	startMediaRPCServer(cfg.RPC.Port, kitexmediahandler.NewHandler(media))
	h := newRouter(media, uploadTasks, chunks, configs.GetServerAddr())
	log.Printf("media-service listening on %s", configs.GetServerAddr())
	h.Spin()
}
