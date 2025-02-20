// models/versions.go
package models

type Version struct {
	ID      uint `gorm:"primaryKey"`
	Type    uint8
	Version string
}
