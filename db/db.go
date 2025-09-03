package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(url string) (*gorm.DB, error) {
	// TODO: change host to docker-compose namespace
	dsn := "postgres://admin:root@localhost/go_otp?sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println("Failed to connect to database:", err)
		return nil, err
	}

	log.Println("Connected to database successfully.")

	return db, nil
}
