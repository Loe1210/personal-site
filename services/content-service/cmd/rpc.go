package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Loe1210/personal-site/internal/xnacos"
	kitexcontentservice "github.com/Loe1210/personal-site/kitex_gen/content/contentservice"
	kitexcontenthandler "github.com/Loe1210/personal-site/services/content-service/internal/handler/rpc"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

type contentServiceRPCConfig struct {
	ServiceName string
	NacosAddr   string
}

func startContentRPCServer(port string, cfg contentServiceRPCConfig, handler *kitexcontenthandler.Handler) {
	addr := resolveTCPAddr(port)
	opts := []server.Option{
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: cfg.ServiceName}),
	}
	registry, err := xnacos.NewRegistry(cfg.NacosAddr)
	if err != nil {
		log.Fatal(err)
	}
	if registry != nil {
		opts = append(opts, server.WithRegistry(registry))
	}

	go func() {
		svr := kitexcontentservice.NewServer(handler, opts...)
		log.Printf("content-service rpc listening on %s", addr.String())
		if err := svr.Run(); err != nil {
			log.Fatal(err)
		}
	}()
}

func contentServiceRPCConfigFromEnv() contentServiceRPCConfig {
	return contentServiceRPCConfig{
		ServiceName: envOrDefault("SERVICE_NAME", "content-service"),
		NacosAddr:   os.Getenv("NACOS_ADDR"),
	}
}

func resolveTCPAddr(port string) net.Addr {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
	return addr
}
func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
