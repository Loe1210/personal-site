package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/gomodule/redigo/redis"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/pkg/xauth"
	"github.com/Loe1210/personal-site/pkg/xotel"
	"github.com/Loe1210/personal-site/services/auth-service/internal/application"
	httpHandler "github.com/Loe1210/personal-site/services/auth-service/internal/handler/http"
	infra "github.com/Loe1210/personal-site/services/auth-service/internal/infra/mysql"
	userrepo "github.com/Loe1210/personal-site/services/auth-service/internal/repository/mysql"
)

var configPath = flag.String("config", "services/auth-service/configs/config.yaml", "auth service config path")

func main() {
	flag.Parse()
	ctx := context.Background()
	cfg, err := configs.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	shutdown, err := xotel.SetupTracerProvider(ctx, "auth-service", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(ctx)
	redisPool := &redis.Pool{
		MaxIdle:     5,
		MaxActive:   20,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.Redis.Addr,
				redis.DialPassword(cfg.Redis.Password),
				redis.DialDatabase(cfg.Redis.DB),
			)
		},
	}
	xauth.UseStore(xauth.NewRedisStore(redisPool, cfg.SessionStore.Prefix))
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
