// models/logs.go
package models

import "gorm.io/gorm"

type LogReportDeploy struct {
	ID                    uint64 `gorm:"primaryKey"`
	Status                uint8
	PackageID             uint
	ClientUniqueID        string
	PreviousLabel         string
	PreviousDeploymentKey string
	CreatedAt             gorm.DeletedAt
}

type LogReportDownload struct {
	ID             uint64 `gorm:"primaryKey"`
	PackageID      uint
	ClientUniqueID string
	CreatedAt      gorm.DeletedAt
}
