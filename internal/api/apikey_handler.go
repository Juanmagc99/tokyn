package api

import (
	"errors"
	"net/http"
	"strings"
	"tokyn/internal/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type APIKeyHandler struct {
	aks service.APIKeyService
}

func NewAPIKeyHandler(aks service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		aks: aks,
	}
}

func (h *APIKeyHandler) Create(c echo.Context) error {
	name := c.FormValue("name")
	if strings.TrimSpace(name) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Please provided a service name for the key")
	}

	tk, err := h.aks.Create(name)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, tk)
}

func (h *APIKeyHandler) Revoke(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing API key ID")
	}

	if err := h.aks.Revoke(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, "API key does not exist or has been revoked")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to revoke API key: "+err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "API key revoked",
		"id":      id,
	})
}

func (h *APIKeyHandler) CheckToken(c echo.Context) error {
	ctx := c.Request().Context()

	token := c.FormValue("token")

	if strings.TrimSpace(token) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing API token")
	}

	ok, err := h.aks.CheckRateLimit(ctx, token)
	if err != nil {
		return err
	}

	if !ok {
		return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
	}

	ak, err := h.aks.CheckToken(ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, "API key does not exist or has been revoked")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error: "+err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"names": ak.Name,
	})
}

func (h *APIKeyHandler) Details(c echo.Context) error {
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing id")
	}

	apk, err := h.aks.Details(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, apk)
}

func (h *APIKeyHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing API key ID")
	}

	if err := h.aks.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, "API key does not exist")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete API key: "+err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "API key deleted",
		"id":      id,
	})
}
