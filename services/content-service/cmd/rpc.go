package main

import (
	"fmt"
	"log"
	"net"

	kitexcontentservice "github.com/Loe1210/personal-site/kitex_gen/content/contentservice"
	kitexcontenthandler "github.com/Loe1210/personal-site/services/content-service/internal/handler/rpc"
	"github.com/cloudwego/kitex/server"
)

func startContentRPCServer(port string, handler *kitexcontenthandler.Handler) {
	go func() {
		addr := resolveTCPAddr(port)
		svr := kitexcontentservice.NewServer(handler, server.WithServiceAddr(addr))
		log.Printf("content-service rpc listening on %s", addr.String())
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
