package models

import (
	"time"

	"gorm.io/gorm"
)

type APIKey struct {
	ID         string         `gorm:"primaryKey"`
	KeyHash    string         `gorm:"not null"`
	Name       string         `gorm:"size:255;not null"`
	CreatedAt  time.Time      `gorm:"not null;autoCreateTime"`
	LastUsedAt *time.Time     `gorm:"default:null"`
	UsageCount int            `gorm:"default:0"`
	Revoked    bool           `gorm:"default:false"`
	RevokedAt  *time.Time     `gorm:"default:null"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
