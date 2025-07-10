package routes

import (
	"log"

	gatewayHandler "github.com/Nick-Spencer-SumUp/test-router/api/handlers/gateway"
	"github.com/Nick-Spencer-SumUp/test-router/internal/config"
	"github.com/labstack/echo/v4"
)

func Gateway(e *echo.Echo) {
	gatewayHandler := gatewayHandler.New()

	// Get all available paths from configuration
	availablePaths := config.GetAllAvailablePaths()

	log.Printf("Registering %d dynamic routes for gateway", len(availablePaths))

	// Register each path with the generic handler
	for _, path := range availablePaths {
		log.Printf("Registering route: %s", path)

		// Register for all HTTP methods - the handler will validate the method
		e.GET(path, gatewayHandler.HandleRequest)
		e.POST(path, gatewayHandler.HandleRequest)
		e.PUT(path, gatewayHandler.HandleRequest)
		e.DELETE(path, gatewayHandler.HandleRequest)
		e.PATCH(path, gatewayHandler.HandleRequest)
	}

	log.Println("Gateway routes registered successfully")
}
