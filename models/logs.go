// models/logs.go
package models

import (
	"time"
)

type LogReportDeploy struct {
	ID                    uint64 `gorm:"primaryKey"`
	Status                uint8
	PackageID             uint
	ClientUniqueID        string
	PreviousLabel         string
	PreviousDeploymentKey string
	CreatedAt             time.Time
}

type LogReportDownload struct {
	ID             uint64 `gorm:"primaryKey"`
	PackageID      uint
	ClientUniqueID string
	CreatedAt      time.Time
}
