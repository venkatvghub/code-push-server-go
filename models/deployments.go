// models/deployments.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Deployment struct {
	ID                      uint `gorm:"primaryKey"`
	AppID                   uint
	Name                    string
	Description             string
	DeploymentKey           string
	LastDeploymentVersionID uint
	LabelID                 uint
	UpdatedAt               time.Time
	CreatedAt               time.Time
	DeletedAt               gorm.DeletedAt
}

type DeploymentHistory struct {
	ID           uint `gorm:"primaryKey"`
	DeploymentID uint
	PackageID    uint
	CreatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

type DeploymentVersion struct {
	ID               uint `gorm:"primaryKey"`
	DeploymentID     uint
	AppVersion       string
	CurrentPackageID uint
	UpdatedAt        time.Time
	CreatedAt        time.Time
	DeletedAt        gorm.DeletedAt
	MinVersion       uint64
	MaxVersion       uint64
}
