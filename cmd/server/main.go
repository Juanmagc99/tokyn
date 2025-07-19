package main

import (
	"tokyn/internal/api"
	"tokyn/internal/models"
	tvalidator "tokyn/pkg/validator"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	sqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.HTTPErrorHandler = api.HandleError
	e.Validator = &tvalidator.CustomValidator{Validator: validator.New()}

	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{TranslateError: true})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	db.AutoMigrate(&models.APIKey{})

	api.NewRouter(e, db)

	e.Logger.Fatal(e.Start("localhost:8080"))
}
