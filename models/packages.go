// models/packages.go
package models

import "gorm.io/gorm"

type Package struct {
	ID                  uint `gorm:"primaryKey"`
	DeploymentVersionID uint
	DeploymentID        uint
	Description         string
	PackageHash         string
	BlobURL             string
	Size                uint
	ManifestBlobURL     string
	ReleaseMethod       string
	Label               string
	OriginalLabel       string
	OriginalDeployment  string
	UpdatedAt           gorm.DeletedAt
	CreatedAt           gorm.DeletedAt
	ReleasedBy          uint64
	IsMandatory         uint8
	IsDisabled          uint8
	Rollout             uint8
	DeletedAt           gorm.DeletedAt
}

type PackageDiff struct {
	ID                     uint `gorm:"primaryKey"`
	PackageID              uint
	DiffAgainstPackageHash string
	DiffBlobURL            string
	DiffSize               uint
	UpdatedAt              gorm.DeletedAt
	CreatedAt              gorm.DeletedAt
	DeletedAt              gorm.DeletedAt
}

type PackageMetrics struct {
	ID         uint `gorm:"primaryKey"`
	PackageID  uint
	Active     uint
	Downloaded uint
	Failed     uint
	Installed  uint
	UpdatedAt  gorm.DeletedAt
	CreatedAt  gorm.DeletedAt
	DeletedAt  gorm.DeletedAt
}
