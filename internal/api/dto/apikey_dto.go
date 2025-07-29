package dto

import (
	"time"
)

type APIKeyCreate struct {
	Name  string `json:"name" validate:"required"`
	Hours int    `json:"hours" validate:"min=0,max=8760"`
}

type APIKeyCreateResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Name  string `json:"name"`
}

type APIKeyDetailsResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	Revoked   bool       `json:"revoked"`
	RevokedAt *time.Time `json:"revoked_at"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type VerifyTokenRequest struct {
	Token string `json:"token"`
}
