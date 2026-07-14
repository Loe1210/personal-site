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
	"github.com/Loe1210/personal-site/services/auth-service/internal/dal/db"
	"github.com/Loe1210/personal-site/services/auth-service/internal/service"
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
	database, err := db.Open(cfg.MySQL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}
	if err := db.Seed(database); err != nil {
		log.Fatal(err)
	}
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
	authService := service.NewAuthService(db.NewUserRepository(database))
	h := server.Default(server.WithHostPorts(configs.GetServerAddr()))
	registerRoutes(h, authService)
	log.Printf("auth-service listening on %s", configs.GetServerAddr())
	h.Spin()
}