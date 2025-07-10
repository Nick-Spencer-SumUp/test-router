package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/Nick-Spencer-SumUp/test-router/internal/config/mappings"
)

type (
	Country       string
	CountryConfig = mappings.ServiceMapping
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

// SelectConfigByPath returns the configuration for a specific country and request path
func SelectConfigByPath(country Country, requestPath string) (CountryConfig, error) {
	if configLoader == nil {
		return CountryConfig{}, fmt.Errorf("configuration not initialized")
	}

	countryConfig, err := configLoader.GetCountryConfig(country)
	if err != nil {
		return CountryConfig{}, fmt.Errorf("failed to get country config: %w", err)
	}

	// Validate that the path exists in the configuration
	if _, err := countryConfig.GetEndpointConfigByPath(requestPath); err != nil {
		return CountryConfig{}, fmt.Errorf("path %s not supported for country %s: %w", requestPath, country, err)
	}

	return countryConfig, nil
}

// SelectConfig returns the configuration for a specific country and route (deprecated - use SelectConfigByPath)
func SelectConfig(country Country, route mappings.Route) (CountryConfig, error) {
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

// GetAllAvailablePaths returns all unique paths available across all countries
func GetAllAvailablePaths() []string {
	if configLoader == nil {
		return nil
	}

	pathSet := make(map[string]bool)
	countries := configLoader.GetAvailableCountries()

	for _, country := range countries {
		countryConfig, err := configLoader.GetCountryConfig(country)
		if err != nil {
			continue
		}

		paths := countryConfig.GetAvailablePaths()
		for _, path := range paths {
			pathSet[path] = true
		}
	}

	// Convert set to slice
	paths := make([]string, 0, len(pathSet))
	for path := range pathSet {
		paths = append(paths, path)
	}

	return paths
}
