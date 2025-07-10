package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/Nick-Spencer-SumUp/test-router/internal/accounts"
	"github.com/Nick-Spencer-SumUp/test-router/internal/config"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	AccountService accounts.Service
}

var (
	BadRequestError          = errors.New("bad request")
	InternalServerError      = errors.New("internal server error")
	CountryNotSupportedError = errors.New("country not supported")
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
	country, err := config.GetCountryFromConfig(countryString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CountryNotSupportedError)
	}

	countryConfig, err := config.SelectConfig(country, config.GetAccountRoute)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError)
	}

	var requestBody accounts.AccountRequest
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError)
	}

	response, err := h.AccountService.GetAccount(countryConfig, requestBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	return h.streamResponse(c, response)
}

func (h *Handler) streamResponse(c echo.Context, response *http.Response) error {
	defer response.Body.Close()

	for key, values := range response.Header {
		for _, value := range values {
			c.Response().Header().Add(key, value)
		}
	}

	c.Response().WriteHeader(response.StatusCode)

	_, err := io.Copy(c.Response().Writer, response.Body)
	if err != nil {
		return err
	}

	return nil
}
