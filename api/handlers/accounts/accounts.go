package handlers

import (
	"net/http"

	"github.com/Nick-Spencer-SumUp/test-router/internal/accounts"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	AccountService accounts.Service
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	InternalServerError       ErrorResponse = ErrorResponse{Error: "internal server error"}
	BadRequestError           ErrorResponse = ErrorResponse{Error: "bad request"}
)

func New(accountService accounts.Service) *Handler {
	return &Handler{
		AccountService: accountService,
	}
}

func (h *Handler) GetAccount(c echo.Context) error {

	var requestBody accounts.AccountRequest
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestError)
	}

	response, err := h.AccountService.GetAccount(c, requestBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError)
	}

	return c.JSON(http.StatusOK, response)
}

