package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/wolbyte/go_otp/db"
	"github.com/wolbyte/go_otp/models"
	"github.com/wolbyte/go_otp/routes"
)

func main() {
	db, err := db.Connect("postgres://admin:root@localhost/go_otp?sslmode=disable")

	// Run AutoMigrate to ensure tables are created
	// TODO: remove
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	routes.Register(router, db)

	router.Run()
}
