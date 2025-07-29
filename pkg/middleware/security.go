package middleware

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func InternalAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		expectedToken := os.Getenv("INTERNAL_API_TOKEN")
		if expectedToken == "" {
			return next(c)
		}

		token := c.Request().Header.Get("X-Internal-Token")
		if token != expectedToken {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
		return next(c)
	}
}
