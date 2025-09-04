package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint      `json:"id" gorm:"uniqueIndex;not null"`
	PhoneNumber  string    `json:"phone_number" gorm:"uniqueIndex;not null"`
	RegisteredAt time.Time `json:"registered_at" gorm:"default:CURRENT_TIMESTAMP"`
	OTP          string
	OTPExpiry    time.Time
}
