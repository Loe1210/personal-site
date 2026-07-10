// @title Personal Site API
// @version 1.0
// @description Personal site backend API for learning Hertz, Kitex and future evolution.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"
    "github.com/hertz-contrib/sessions"
    "github.com/hertz-contrib/sessions/cookie"

	"github.com/Loe1210/personal-site/biz/dal/db"
	"github.com/Loe1210/personal-site/biz/router"
)


func main() {
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	h := server.Default()
	store := cookie.NewStore([]byte("personal-site-session-secret"))
    h.Use(sessions.New("personal_site_session", store))


	h.StaticFile("/swagger.json", "./docs/swagger.json")
	h.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
	h.Static("/static", "./static")
	
	router.Register(h)

	h.Spin()
}