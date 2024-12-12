package models

import (
	"time"

	"gorm.io/gorm"
)

// System User model
type User struct {
	ID           uint           `gorm:"primaryKey"`
	Name         string         `gorm:"size:250;"`
	Username     string         `gorm:"uniqueIndex;size:100;not null"`
	Email        string         `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash []byte         `gorm:"not null"`
	IsStaff      bool           `gorm:"default:false"` // For staff members
	IsSuperuser  bool           `gorm:"default:false"` // For admins
	IsActive     bool           `gorm:"default:true"`  // Can be banned or active
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoCreateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
