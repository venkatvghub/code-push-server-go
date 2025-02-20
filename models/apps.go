// models/apps.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type App struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	UID           uint64
	OS            uint8
	Platform      uint8
	IsUseDiffText uint8
	UpdatedAt     time.Time
	CreatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
