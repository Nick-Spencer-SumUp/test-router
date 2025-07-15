package accounts

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Nick-Spencer-SumUp/test-router/internal/config"
	"github.com/labstack/echo/v4"
)

type (
	AccountRequest struct {
		Mid string `json:"mid"`
	}

	AccountResponse struct {
		AccountUID string `json:"account_uid"`
		Status     string `json:"status"`
		DeletedAt  string `json:"deleted_at"`
		IsLocked   bool   `json:"is_locked"`
	}
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s Service) GetAccount(ctx echo.Context, accountRequest AccountRequest) (*AccountResponse, error) {
	countryConfig := ctx.Get("countryConfig").(config.Service)

	proxyPath := countryConfig.BaseUrl + ctx.Request().URL.Path
	method := ctx.Request().Method

	requestBody, err := json.Marshal(accountRequest)
	if err != nil {
		return nil, err
	}

	response, err := s.doRequest(method, proxyPath, requestBody)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	// validate the response
	var accountResponse AccountResponse
	err = json.NewDecoder(response.Body).Decode(&accountResponse)
	if err != nil {
		return nil, err
	}

	return &accountResponse, nil
}

func (s Service) doRequest(method string, path string, requestBody []byte) (*http.Response, error) {

	request, err := http.NewRequest(method, path, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(request)
}
