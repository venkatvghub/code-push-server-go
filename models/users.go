// models/users.go
package models

import "gorm.io/gorm"

type UserToken struct {
	ID          uint64 `gorm:"primaryKey"`
	UID         uint64
	Name        string
	Tokens      string
	CreatedBy   string
	Description string
	IsSession   uint8
	ExpiresAt   gorm.DeletedAt
	CreatedAt   gorm.DeletedAt
	DeletedAt   gorm.DeletedAt
}

type User struct {
	ID        uint64 `gorm:"primaryKey"`
	Username  string
	Password  string
	Email     string
	Identical string
	AckCode   string
	UpdatedAt gorm.DeletedAt
	CreatedAt gorm.DeletedAt
}
