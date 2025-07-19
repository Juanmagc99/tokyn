package api

import (
	"tokyn/internal/repository"
	"tokyn/internal/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewRouter(e *echo.Echo, db *gorm.DB) {

	apiKeyRepo := repository.NewAPIKeyRepository(db)
	apiKeyService := service.NewAPIKeyService(*apiKeyRepo)
	apiKeyHandler := NewAPIKeyHandler(*apiKeyService)

	e.POST("/api/token", apiKeyHandler.Create)
}
