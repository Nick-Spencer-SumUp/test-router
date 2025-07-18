package accounts

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Nick-Spencer-SumUp/test-router/internal/config/countries"
	"github.com/Nick-Spencer-SumUp/test-router/internal/config/mappings"
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

func (s Service) GetAccount(cfg countries.CountryConfig, accountRequest AccountRequest) (*http.Response, error) {
	requestBody, err := json.Marshal(accountRequest)
	if err != nil {
		return nil, err
	}

	endpointConfig, err := cfg.GetEndpointConfig(mappings.GetAccountRoute)
	if err != nil {
		return nil, err
	}

	response, err := s.doRequest(endpointConfig.GetMethod(), endpointConfig, requestBody)
	if err != nil {
		return nil, err
	}

	// Return the raw response without decoding it for streaming
	return response, nil
}

func (s Service) doRequest(method string, endpoint mappings.EndpointConfig, requestBody []byte) (*http.Response, error) {
	request, err := http.NewRequest(method, endpoint.BaseURL+endpoint.Endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(request)
}
