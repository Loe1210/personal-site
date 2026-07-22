package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/internal/xsafe"
	"github.com/Loe1210/personal-site/services/content-service/internal/dal/db"
	kitexcontenthandler "github.com/Loe1210/personal-site/services/content-service/internal/handler/rpc"
	"github.com/Loe1210/personal-site/services/content-service/internal/service"
	"github.com/Loe1210/personal-site/services/content-service/pkg/xotel"
)

var configPath = flag.String("config", "services/content-service/configs/config.yaml", "content service config path")

func main() {
	flag.Parse()
	xsafe.InstallGoPoolPanicHandler()
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
	articles := service.NewArticleServiceWithCaches(repo, service.NewLocalArticleCache(), service.NewRedisArticleCache(redisPool, "content:article:", 10*time.Minute))
	categories := service.NewCategoryService(repo)
	tags := service.NewTagService(repo)
	startContentRPCServer(cfg.RPC.Port, contentServiceRPCConfigFromEnv(), kitexcontenthandler.NewHandler(articles))
	h := newRouter(articles, categories, tags, configs.GetServerAddr())
	log.Printf("content-service listening on %s", configs.GetServerAddr())
	h.Spin()
}
