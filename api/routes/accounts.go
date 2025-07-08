package routes

import (
	"github.com/labstack/echo/v4"
	accountsHandler "github.com/Nick-Spencer-SumUp/test-router/api/handlers/accounts"
	"github.com/Nick-Spencer-SumUp/test-router/internal/accounts"
)

func Accounts(e *echo.Group) {
	accountsService := accounts.New()
	accountsHandler := accountsHandler.New(*accountsService)

	e.GET("", accountsHandler.GetAccount)
}
