package db

import (
	"fmt"
	"log"

	"github.com/wolbyte/go_otp/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(url string) (*gorm.DB, error) {
	dbUsername := utils.GetenvDefault("DATABASE_USERNAME", "admin")
	dbPassword := utils.GetenvDefault("DATABASE_PASSWORD", "root")
	dbName := utils.GetenvDefault("DATABASE_NAME", "go_otp")
	dbHost := utils.GetenvDefault("DATABASE_HOST", "db")

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUsername, dbPassword, dbHost, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println("Failed to connect to database:", err)
		return nil, err
	}

	log.Println("Connected to database successfully.")

	return db, nil
}
