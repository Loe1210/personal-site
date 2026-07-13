package main

import (
	"flag"
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/services/gateway/internal/router"
)

var configPath = flag.String("config", "configs/config.yaml", "gateway config path")

func main() {
	flag.Parse()
	_, err := configs.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	h := server.Default(server.WithHostPorts(configs.GetServerAddr()))
	deps := router.Dependencies{AuthServiceName: "auth-service", BFFServiceName: "web-bff"}
	if err := router.RegisterRoutes(h, deps); err != nil {
		log.Fatal(err)
	}
	log.Printf("gateway listening on %s", configs.GetServerAddr())
	h.Spin()
}
