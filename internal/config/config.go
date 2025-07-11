package config

import (
	"fmt"
	"log"
	"sync"
)

// Core configuration types
type (
	Route  string
	Method string

	ServiceMapping struct {
		Endpoints map[Route]Endpoint
		BaseURL   string
	}

	Endpoint struct {
		Method Method
		URI    string
	}

	EndpointConfig struct {
		Endpoint string
		BaseURL  string
		Method   Method
	}
)

// Route and Method constants
const (
	GetAccountRoute Route = "GetAccount"

	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
)

var MethodMap = map[Method]string{
	GET:    "GET",
	POST:   "POST",
	PUT:    "PUT",
	DELETE: "DELETE",
}

// ServiceMapping methods
func (c ServiceMapping) GetEndpointConfig(route Route) (EndpointConfig, error) {
	endpoint, ok := c.Endpoints[route]
	if !ok {
		return EndpointConfig{}, fmt.Errorf("route not supported")
	}
	return EndpointConfig{
		Endpoint: endpoint.URI,
		BaseURL:  c.BaseURL,
		Method:   endpoint.Method,
	}, nil
}

func (e EndpointConfig) GetMethod() string {
	return MethodMap[e.Method]
}

// Configuration management types
type (
	Country       string
	CountryConfig = ServiceMapping
	RoutesConfig  = map[Country]CountryConfig
)

var (
	configLoader *ConfigLoader
	initOnce     sync.Once
)

// InitConfig initializes the configuration loader
func InitConfig(configDir string) error {
	var err error
	initOnce.Do(func() {
		configLoader = NewConfigLoader(configDir)
		err = configLoader.LoadConfig()
		if err != nil {
			return
		}
		err = configLoader.ValidateConfig()
		if err != nil {
			return
		}
		log.Println("Configuration loaded successfully")
	})
	return err
}

// GetCountryFromConfig retrieves a country from the YAML configuration
func GetCountryFromConfig(countryString string) (Country, error) {
	if configLoader == nil {
		return "", fmt.Errorf("configuration not initialized")
	}

	countries := configLoader.GetAvailableCountries()
	for _, country := range countries {
		if string(country) == countryString {
			return country, nil
		}
	}

	return "", fmt.Errorf("country %s not found in configuration", countryString)
}

// SelectConfig returns the configuration for a specific country and route
func SelectConfig(country string, route Route) (CountryConfig, error) {
	if configLoader == nil {
		return CountryConfig{}, fmt.Errorf("configuration not initialized")
	}

	countryConfig, err := configLoader.GetCountryConfig(country)
	if err != nil {
		return CountryConfig{}, fmt.Errorf("failed to get country config: %w", err)
	}

	// Validate that the route exists in the configuration
	if _, err := countryConfig.GetEndpointConfig(route); err != nil {
		return CountryConfig{}, fmt.Errorf("route %s not supported for country %s: %w", route, country, err)
	}

	return countryConfig, nil
}

// GetAvailableCountries returns all available countries from the configuration
func GetAvailableCountries() []Country {
	if configLoader == nil {
		return nil
	}
	return configLoader.GetAvailableCountries()
}

// ReloadConfig reloads the configuration from the file
func ReloadConfig() error {
	if configLoader == nil {
		return fmt.Errorf("configuration not initialized")
	}
	return configLoader.ReloadConfig()
}
