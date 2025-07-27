package models

import (
	"time"
)

type APIKey struct {
	ID        string     `gorm:"primaryKey"`
	KeyHash   string     `gorm:"not null"`
	Name      string     `gorm:"size:255;not null"`
	CreatedAt time.Time  `gorm:"not null;autoCreateTime"`
	Revoked   bool       `gorm:"default:false"`
	RevokedAt *time.Time `gorm:"default:null"`
}

type APIKeyRedis struct {
	ID      string `redis:"primaryKey"`
	KeyHash string `redis:"keyHash"`
	Name    string `redis:"name"`
}
