package api

import (
	"tokyn/internal/repository"
	"tokyn/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(e *echo.Echo, db *gorm.DB, rclient *redis.Client) {

	apiKeyRepo := repository.NewAPIKeyRepository(db, rclient)
	apiKeyService := service.NewAPIKeyService(*apiKeyRepo)
	apiKeyHandler := NewAPIKeyHandler(*apiKeyService)

	e.POST("/apikeys", apiKeyHandler.Create)
	e.DELETE("/apikeys/:id", apiKeyHandler.Revoke)
	e.GET("/apikeys/verify/:token", apiKeyHandler.CheckToken)
	e.GET("/apikeys/:id/details", apiKeyHandler.Details)
}
