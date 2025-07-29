package main

import (
	"tokyn/config"
	"tokyn/internal/api"
	"tokyn/internal/models"
	"tokyn/internal/storage"
	tvalidator "tokyn/pkg/validator"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	cfg := config.NewConfig()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.HTTPErrorHandler = api.HandleError
	e.Validator = &tvalidator.CustomValidator{Validator: validator.New()}

	db, err := storage.NewSQLDB(cfg.SQLiteDB)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	rclient, err := storage.NewRedisDB(cfg.RedisAddr, cfg.RedisPass)
	if err != nil {
		panic("failed to connect redis: " + err.Error())
	}

	db.AutoMigrate(&models.APIKey{})

	api.NewRouter(e, db, rclient)

	e.Logger.Fatal(e.Start(cfg.AppAddr))
}
