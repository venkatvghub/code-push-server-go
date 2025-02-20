// models/collaborators.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Collaborator struct {
	ID        uint64 `gorm:"primaryKey"`
	AppID     uint
	UID       uint64
	Roles     string
	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}
