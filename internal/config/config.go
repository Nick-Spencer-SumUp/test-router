package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/Nick-Spencer-SumUp/test-router/internal/config/countries"
	"github.com/Nick-Spencer-SumUp/test-router/internal/config/mappings"
)

type RoutesConfig map[countries.Country]countries.CountryConfig

var (
	configLoader *ConfigLoader
	initOnce     sync.Once
)

// InitConfig initializes the configuration loader
func InitConfig(configPath string) error {
	var err error
	initOnce.Do(func() {
		configLoader = NewConfigLoader(configPath)
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

// SelectConfig returns the configuration for a specific country and route
func SelectConfig(country countries.Country, route mappings.Route) (countries.CountryConfig, error) {
	if configLoader == nil {
		return countries.CountryConfig{}, fmt.Errorf("configuration not initialized")
	}

	countryConfig, err := configLoader.GetCountryConfig(country)
	if err != nil {
		return countries.CountryConfig{}, fmt.Errorf("failed to get country config: %w", err)
	}

	// Validate that the route exists in the configuration
	if _, err := countryConfig.GetEndpointConfig(route); err != nil {
		return countries.CountryConfig{}, fmt.Errorf("route %s not supported for country %s: %w", route, country, err)
	}

	return countryConfig, nil
}

// GetAvailableCountries returns all available countries from the configuration
func GetAvailableCountries() []countries.Country {
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
