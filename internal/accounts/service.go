package accounts

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sumup/test-router/internal/config"
)

type (
	AccountRequest struct {
		Mid string `json:"mid"`
	}

	AccountInfo struct {
		AccountNumber  string `json:"account_number"`
		AccountName    string `json:"account_name"`
		AccountType    string `json:"account_type"`
		AccountBalance string `json:"account_balance"`
	}
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s Service) GetAccount(config config.EndpointConfig, accountRequest AccountRequest) (AccountInfo, error) {
	requestBody, err := json.Marshal(accountRequest)
	if err != nil {
		return AccountInfo{}, err
	}

	response, err := s.doRequest("GET", config, requestBody)
	if err != nil {
		return AccountInfo{}, err
	}

	var accountInfo AccountInfo
	err = json.NewDecoder(response.Body).Decode(&accountInfo)
	if err != nil {
		return AccountInfo{}, err
	}

	return accountInfo, nil
}

func (s Service) doRequest(method string, config config.EndpointConfig, requestBody []byte) (*http.Response, error) {
	request, err := http.NewRequest(method, config.BaseURL+config.Endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(request)
}
