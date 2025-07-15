package main

import (
	"log"

	accountsHandler "github.com/Nick-Spencer-SumUp/test-router/api/handlers/accounts"
	accountsService "github.com/Nick-Spencer-SumUp/test-router/internal/accounts"
	authMiddleware "github.com/Nick-Spencer-SumUp/test-router/internal/middleware"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize configuration
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize echo server
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(authMiddleware.SetConfigFromToken)

	accountsService := accountsService.New()
	accountsHandler := accountsHandler.New(*accountsService)

	// Setup routes
	e.GET("/jokes/random", accountsHandler.GetAccount)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
