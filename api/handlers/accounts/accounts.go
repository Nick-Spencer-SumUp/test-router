package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sumup/test-router/internal/accounts"
	"github.com/sumup/test-router/internal/config"
)

type Handler struct {
	AccountService accounts.Service
}

var (
	BadRequestError = errors.New("bad request")
	InternalServerError = errors.New("internal server error")
)

func New(accountService accounts.Service) *Handler {
	return &Handler{
		AccountService: accountService,
	}
}

func (h *Handler) GetAccount(c echo.Context) error {
	// TODO: get locale from request context, likely from token claims or internal api call

	// TODO: decide, should config be selected in middlware and passed to handler/context?
	countryString := c.Request().Header.Get("country")
	country := config.GetCountry(countryString)

	routeConfig, err := config.SelectConfig(country, config.GetAccountRoute)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError)
	}

	var requestBody accounts.AccountRequest
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError)
	}

	response, err := h.AccountService.GetAccount(routeConfig, requestBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	return c.JSON(http.StatusOK, response)
}
