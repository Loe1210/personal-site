package main

import (
	"fmt"
	"log"
	"net"

	kitexauthservice "github.com/Loe1210/personal-site/kitex_gen/auth/authservice"
	kitexauthhandler "github.com/Loe1210/personal-site/services/auth-service/internal/handler/rpc"
	"github.com/cloudwego/kitex/server"
)

func startAuthRPCServer(port string, handler *kitexauthhandler.Handler) {
	go func() {
		addr := resolveTCPAddr(port)
		svr := kitexauthservice.NewServer(handler, server.WithServiceAddr(addr))
		log.Printf("auth-service rpc listening on %s", addr.String())
		if err := svr.Run(); err != nil {
			log.Fatal(err)
		}
	}()
}

func resolveTCPAddr(port string) net.Addr {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
	return addr
}
