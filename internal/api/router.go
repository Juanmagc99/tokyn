package api

import (
	"tokyn/internal/repository"
	"tokyn/internal/service"
	"tokyn/pkg/middleware"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(e *echo.Echo, db *gorm.DB, rclient *redis.Client) {

	apiKeyRepo := repository.NewAPIKeyRepository(db, rclient)
	apiKeyService := service.NewAPIKeyService(*apiKeyRepo)
	apiKeyHandler := NewAPIKeyHandler(*apiKeyService)

	g := e.Group("/apikeys", middleware.InternalAuth)
	g.POST("", apiKeyHandler.Create)
	g.GET("/:id", apiKeyHandler.Details)
	g.PATCH("/:id/revocation", apiKeyHandler.Revoke)
	g.DELETE("/:id", apiKeyHandler.Delete)
	g.GET("/verification", apiKeyHandler.CheckToken)

}
