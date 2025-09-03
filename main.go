package main

import (
	"github.com/gin-gonic/gin"

	"github.com/wolbyte/go_otp/routes"
)

func main() {
	router := gin.New()
	router.Use(gin.Recovery())
	routes.Register(router)

	router.Run()
}
