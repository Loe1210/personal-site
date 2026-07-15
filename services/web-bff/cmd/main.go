package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/services/web-bff/internal/assembler"
	httpHandler "github.com/Loe1210/personal-site/services/web-bff/internal/handler/http"
	"github.com/Loe1210/personal-site/services/web-bff/pkg/xotel"
)

var (
	configPath        = flag.String("config", "services/web-bff/configs/config.yaml", "web bff config path")
	contentServiceURL = flag.String("content-service-url", "http://127.0.0.1:9003", "content service base url")
)

func main() {
	flag.Parse()
	ctx := context.Background()
	_, err := configs.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	shutdown, err := xotel.SetupTracerProvider(ctx, "web-bff", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(ctx)
	articleAssembler := assembler.NewArticlePageAssembler(assembler.NewHTTPContentClient(*contentServiceURL))
	handler := httpHandler.NewHandler(articleAssembler)
	h := server.Default(server.WithHostPorts(configs.GetServerAddr()))
	handler.RegisterRoutes(h)
	log.Printf("web-bff listening on %s", configs.GetServerAddr())
	h.Spin()
}
