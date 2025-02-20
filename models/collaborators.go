// models/collaborators.go
package models

import "gorm.io/gorm"

type Collaborator struct {
	ID        uint64 `gorm:"primaryKey"`
	AppID     uint
	UID       uint64
	Roles     string
	UpdatedAt gorm.DeletedAt
	CreatedAt gorm.DeletedAt
	DeletedAt gorm.DeletedAt
}
