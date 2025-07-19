package api

import (
	"net/http"
	"strings"
	"tokyn/internal/service"

	"github.com/labstack/echo/v4"
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

	return c.JSON(http.StatusCreated, echo.Map{
		"api_key": tk,
	})
}
