package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nick-Spencer-SumUp/test-router/internal/config"
	"github.com/labstack/echo/v4"
)

type (
	AccountRequest struct {
		Mid string `json:"mid"`
	}
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s Service) GetAccount(ctx echo.Context, accountRequest AccountRequest) (*http.Response, error) {
	countryConfig := ctx.Get("countryConfig").(config.CountryConfig)

	endpointConfig, err := countryConfig.GetEndpointConfig(config.GetAccountRoute)
	if err != nil {
		return nil, err
	}

	requestBody, err := json.Marshal(accountRequest)
	if err != nil {
		return nil, err
	}

	// No path parameters needed for GetAccount
	response, err := s.doRequest(endpointConfig.GetMethod(), endpointConfig, requestBody, nil)
	if err != nil {
		return nil, err
	}

	// Return the raw response without decoding it for streaming
	return response, nil
}

// doRequest handles HTTP requests with optional path parameter substitution using fmt.Sprintf
func (s Service) doRequest(method string, endpoint config.EndpointConfig, requestBody []byte, pathParams []interface{}) (*http.Response, error) {
	// Build the URL with path parameters if provided
	var finalURL string
	if len(pathParams) > 0 {
		finalURL = endpoint.BaseURL + fmt.Sprintf(endpoint.Endpoint, pathParams...)
	} else {
		finalURL = endpoint.BaseURL + endpoint.Endpoint
	}

	request, err := http.NewRequest(method, finalURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(request)
}
