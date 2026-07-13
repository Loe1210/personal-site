package main

import (
	"flag"
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/services/auth-service/internal/application"
	httpHandler "github.com/Loe1210/personal-site/services/auth-service/internal/handler/http"
	infra "github.com/Loe1210/personal-site/services/auth-service/internal/infra/mysql"
	userrepo "github.com/Loe1210/personal-site/services/auth-service/internal/repository/mysql"
)

var configPath = flag.String("config", "services/auth-service/configs/config.yaml", "auth service config path")

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
	service := application.NewAuthService(userrepo.NewUserRepository(database))
	handler := httpHandler.NewHandler(service)
	h := server.Default(server.WithHostPorts(configs.GetServerAddr()))
	handler.RegisterRoutes(h)
	log.Printf("auth-service listening on %s", configs.GetServerAddr())
	h.Spin()
}
