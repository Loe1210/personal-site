package main

import (
	"flag"
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/services/media-service/internal/application"
	httpHandler "github.com/Loe1210/personal-site/services/media-service/internal/handler/http"
	infra "github.com/Loe1210/personal-site/services/media-service/internal/infra/mysql"
	"github.com/Loe1210/personal-site/services/media-service/internal/infra/storage"
	filerepo "github.com/Loe1210/personal-site/services/media-service/internal/repository/mysql"
)

var configPath = flag.String("config", "services/media-service/configs/config.yaml", "media service config path")

func main() {
	flag.Parse()
	cfg, err := configs.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	database, err := infra.Open(cfg.MySQL)
	if err != nil {
		log.Fatal(err)
	}
	store := storage.NewLocalStorage(cfg.Upload.RootDir, cfg.Upload.PublicBasePath)
	service := application.NewMediaService(store, filerepo.NewFileRepository(database))
	handler := httpHandler.NewHandler(service)
	h := server.Default(server.WithHostPorts(configs.GetServerAddr()))
	handler.RegisterRoutes(h)
	log.Printf("media-service listening on %s", configs.GetServerAddr())
	h.Spin()
}
