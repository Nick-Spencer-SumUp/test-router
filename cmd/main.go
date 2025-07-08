package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/Nick-Spencer-SumUp/test-router/api/routes"
)

func main() {

	// Initialize echo server
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Setup routes
	accountGroup := e.Group("/accounts")
	routes.Accounts(accountGroup)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
