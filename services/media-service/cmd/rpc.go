package main

import (
	"fmt"
	"log"
	"net"

	kitexmediaservice "github.com/Loe1210/personal-site/kitex_gen/media/mediaservice"
	kitexmediahandler "github.com/Loe1210/personal-site/services/media-service/internal/handler/rpc"
	"github.com/cloudwego/kitex/server"
)

func startMediaRPCServer(port string, handler *kitexmediahandler.Handler) {
	go func() {
		addr := resolveTCPAddr(port)
		svr := kitexmediaservice.NewServer(handler, server.WithServiceAddr(addr))
		log.Printf("media-service rpc listening on %s", addr.String())
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
