package admin

import (
	"net/http"

	"github.com/Nick-Spencer-SumUp/test-router/internal/config"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

// ReloadConfig provides a hot-reload endpoint for configuration
func (h *Handler) ReloadConfig(c echo.Context) error {
	if err := config.ReloadConfig(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to reload configuration",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Configuration reloaded successfully",
	})
}

// GetAvailableCountries returns all available countries from the configuration
func (h *Handler) GetAvailableCountries(c echo.Context) error {
	countries := config.GetAvailableCountries()

	countryList := make([]string, len(countries))
	for i, country := range countries {
		countryList[i] = string(country)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"countries": countryList,
		"count":     len(countryList),
	})
}

// HealthCheck provides a health check endpoint
func (h *Handler) HealthCheck(c echo.Context) error {
	countries := config.GetAvailableCountries()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "healthy",
		"config": map[string]interface{}{
			"available_countries": len(countries),
			"config_loaded":       countries != nil,
		},
	})
}
