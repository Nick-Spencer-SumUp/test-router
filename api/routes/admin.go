package routes

import (
	adminHandler "github.com/Nick-Spencer-SumUp/test-router/api/handlers/admin"
	"github.com/labstack/echo/v4"
)

func Admin(e *echo.Group) {
	adminHandler := adminHandler.New()

	// Hot-reload configuration
	e.POST("/reload-config", adminHandler.ReloadConfig)

	// Get available countries
	e.GET("/countries", adminHandler.GetAvailableCountries)

	// Health check
	e.GET("/health", adminHandler.HealthCheck)
}
