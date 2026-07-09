// @title Personal Site API
// @version 1.0
// @description Personal site backend API for learning Hertz, Kitex and future evolution.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Loe1210/personal-site/biz/router"
)


func main() {
	h := server.Default()

	h.StaticFile("/swagger.json", "./docs/swagger.json")
	h.StaticFile("/swagger.yaml", "./docs/swagger.yaml")

	router.Register(h)

	h.Spin()
}