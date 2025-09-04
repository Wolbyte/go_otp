package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	PhoneNumber  string    `json:"phone_number" gorm:"uniqueIndex;not null"`
	RegisteredAt time.Time `json:"registered_at" gorm:"default:CURRENT_TIMESTAMP"`
}
