// models/deployments.go
package models

import "gorm.io/gorm"

type Deployment struct {
	ID                      uint `gorm:"primaryKey"`
	AppID                   uint
	Name                    string
	Description             string
	DeploymentKey           string
	LastDeploymentVersionID uint
	LabelID                 uint
	UpdatedAt               gorm.DeletedAt
	CreatedAt               gorm.DeletedAt
	DeletedAt               gorm.DeletedAt
}

type DeploymentHistory struct {
	ID           uint `gorm:"primaryKey"`
	DeploymentID uint
	PackageID    uint
	CreatedAt    gorm.DeletedAt
	DeletedAt    gorm.DeletedAt
}

type DeploymentVersion struct {
	ID               uint `gorm:"primaryKey"`
	DeploymentID     uint
	AppVersion       string
	CurrentPackageID uint
	UpdatedAt        gorm.DeletedAt
	CreatedAt        gorm.DeletedAt
	DeletedAt        gorm.DeletedAt
	MinVersion       uint64
	MaxVersion       uint64
}
