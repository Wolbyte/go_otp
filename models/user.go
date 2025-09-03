package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"uniqueIndex;not null"`
	PhoneNumber string `json:"phoneNumber" gorm:"uniqueIndex;not null"`
	OTP         string
	OTPExpiry   time.Time
}
