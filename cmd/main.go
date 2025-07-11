package main

import (
	"log"
	"os"

	"github.com/Nick-Spencer-SumUp/test-router/api/routes"
	"github.com/Nick-Spencer-SumUp/test-router/internal/config"
	authMiddleware "github.com/Nick-Spencer-SumUp/test-router/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize configuration
	configDir := os.Getenv("CONFIG_DIR")
	if configDir == "" {
		configDir = "internal/config"
	}

	if err := config.InitConfig(configDir); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Initialize echo server
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(authMiddleware.GetCountryFromToken)

	// Setup routes
	accountGroup := e.Group("/accounts")
	routes.Accounts(accountGroup)

	// Setup admin routes
	adminGroup := e.Group("/admin")
	routes.Admin(adminGroup)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
