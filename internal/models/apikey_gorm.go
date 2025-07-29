package models

import (
	"errors"
	"time"
)

var ErrAPIKeyInvalid = errors.New("API key is revoked or expired")

type APIKey struct {
	ID        string     `gorm:"primaryKey"`
	KeyHash   string     `gorm:"not null"`
	Name      string     `gorm:"size:255;not null"`
	CreatedAt time.Time  `gorm:"not null;autoCreateTime"`
	Revoked   bool       `gorm:"default:false"`
	RevokedAt *time.Time `gorm:"default:null"`
	ExpiresAt *time.Time `gorm:"default:null"`
}

func (a *APIKey) IsValid() bool {
	if a.ExpiresAt != nil && a.ExpiresAt.After(time.Now()) {
		return true
	}

	return false
}
