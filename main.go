// @title			GO OTP
// @version		1.0
// @description	A simple OTP implementation written in go
// @host			localhost:8080
// @BasePath		/api/v1
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/wolbyte/go_otp/db"
	"github.com/wolbyte/go_otp/models"
	"github.com/wolbyte/go_otp/routes"
	"github.com/wolbyte/go_otp/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, falling back to system environment")
	}

	gin.SetMode(utils.GetenvDefault("GIN_MODE", "release"))

	db, err := db.Connect("postgres://admin:root@localhost/go_otp?sslmode=disable")

	// Run AutoMigrate to ensure tables are created
	// golang-migrate can be used for a more robust approach, but for the sake of simplicity we just use gorm's built-in method
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
