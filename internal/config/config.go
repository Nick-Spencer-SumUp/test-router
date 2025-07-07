package config

import "fmt"

type (
	Country string
	Route   string
	Method  string

	EndpointConfig struct {
		Endpoint string
		BaseURL  string
	}

	CountryConfig struct {
		Endpoints map[Route]Endpoint
		BaseURL   string
	}

	Endpoint struct {
		Method Method
		URI    string
	}
)

type RoutesConfig map[Country]CountryConfig

const (
	US Country = "US"
	DE Country = "DE"

	GetAccountRoute Route = "GetAccount"

	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
)

var (
	RouterConfigs = RoutesConfig{
		US: usConfig,
		DE: deConfig,
	}

	CountryMap = map[string]Country{
		"US": US,
		"DE": DE,
	}
)

func SelectConfig(country Country, route Route) (EndpointConfig, error) {
	countryConfig, ok := RouterConfigs[country]
	if !ok {
		return EndpointConfig{}, fmt.Errorf("country not supported")
	}

	endpoint, ok := countryConfig.Endpoints[route]
	if !ok {
		return EndpointConfig{}, fmt.Errorf("route not supported")
	}

	return EndpointConfig{
		Endpoint: endpoint.URI,
		BaseURL:  countryConfig.BaseURL,
	}, nil
}

func GetCountry(country string) Country {
	return CountryMap[country]
}
