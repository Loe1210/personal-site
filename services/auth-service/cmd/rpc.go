package main

import (
	"fmt"
	"log"
	"net"
	"os"

	kitexauthservice "github.com/Loe1210/personal-site/kitex_gen/auth/authservice"
	kitexauthhandler "github.com/Loe1210/personal-site/services/auth-service/internal/handler/rpc"
	"github.com/Loe1210/personal-site/services/auth-service/pkg/xnacos"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

type authServiceRPCConfig struct {
	ServiceName string
	NacosAddr   string
}

func startAuthRPCServer(port string, cfg authServiceRPCConfig, handler *kitexauthhandler.Handler) {
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
		svr := kitexauthservice.NewServer(handler, opts...)
		log.Printf("auth-service rpc listening on %s", addr.String())
		if err := svr.Run(); err != nil {
			log.Fatal(err)
		}
	}()
}

func authServiceRPCConfigFromEnv() authServiceRPCConfig {
	return authServiceRPCConfig{
		ServiceName: envOrDefault("SERVICE_NAME", "auth-service"),
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
