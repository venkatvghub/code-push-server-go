// models/users.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type UserToken struct {
	ID          uint64 `gorm:"primaryKey"`
	UID         uint64
	Name        string
	Tokens      string
	CreatedBy   string
	Description string
	IsSession   uint8
	ExpiresAt   gorm.DeletedAt
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type User struct {
	ID        uint64 `gorm:"primaryKey"`
	Username  string
	Password  string
	Email     string
	Identical string
	AckCode   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
