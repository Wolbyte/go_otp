package models

import (
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	PhoneNumber  string    `json:"phone_number" gorm:"uniqueIndex;not null"`
	RegisteredAt time.Time `json:"registered_at" gorm:"default:CURRENT_TIMESTAMP"`
}
