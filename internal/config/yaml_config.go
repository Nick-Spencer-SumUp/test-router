package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Nick-Spencer-SumUp/test-router/internal/config/countries"
	"github.com/Nick-Spencer-SumUp/test-router/internal/config/mappings"
	"gopkg.in/yaml.v3"
)

// YAMLConfig represents the structure of the YAML configuration file
type YAMLConfig struct {
	Services     map[string]YAMLService     `yaml:"services"`
	Countries    map[string]YAMLCountry     `yaml:"countries"`
	Environments map[string]YAMLEnvironment `yaml:"environments"`
}

type YAMLService struct {
	BaseURL   string                  `yaml:"base_url"`
	Endpoints map[string]YAMLEndpoint `yaml:"endpoints"`
}

type YAMLEndpoint struct {
	Method string `yaml:"method"`
	URI    string `yaml:"uri"`
}

type YAMLCountry struct {
	Service  string   `yaml:"service"`
	Region   string   `yaml:"region"`
	Features []string `yaml:"features"`
}

type YAMLEnvironment struct {
	Services map[string]YAMLService `yaml:"services"`
}

// ConfigLoader handles loading and parsing YAML configuration
type ConfigLoader struct {
	configPath string
	config     *YAMLConfig
	mu         sync.RWMutex
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader(configPath string) *ConfigLoader {
	return &ConfigLoader{
		configPath: configPath,
	}
}

// LoadConfig loads the configuration from the YAML file
func (cl *ConfigLoader) LoadConfig() error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	// Get environment
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	// Read the config file
	configFile := cl.configPath
	if configFile == "" {
		configFile = "config/routing.yaml"
	}

	// Get absolute path
	absPath, err := filepath.Abs(configFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", absPath, err)
	}

	// Parse YAML
	var config YAMLConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse YAML config: %w", err)
	}

	// Apply environment-specific overrides
	if envConfig, exists := config.Environments[env]; exists {
		for serviceName, envService := range envConfig.Services {
			if service, exists := config.Services[serviceName]; exists {
				// Override base URL if specified
				if envService.BaseURL != "" {
					service.BaseURL = envService.BaseURL
					config.Services[serviceName] = service
				}
				// Could add more overrides here (endpoints, etc.)
			}
		}
	}

	cl.config = &config
	return nil
}

// GetCountryConfig returns the configuration for a specific country
func (cl *ConfigLoader) GetCountryConfig(country countries.Country) (countries.CountryConfig, error) {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	if cl.config == nil {
		return countries.CountryConfig{}, fmt.Errorf("configuration not loaded")
	}

	countryStr := string(country)
	countryConfig, exists := cl.config.Countries[countryStr]
	if !exists {
		return countries.CountryConfig{}, fmt.Errorf("country %s not found", countryStr)
	}

	serviceConfig, exists := cl.config.Services[countryConfig.Service]
	if !exists {
		return countries.CountryConfig{}, fmt.Errorf("service %s not found for country %s", countryConfig.Service, countryStr)
	}

	// Convert YAML config to ServiceMapping
	serviceMapping := mappings.ServiceMapping{
		BaseURL:   serviceConfig.BaseURL,
		Endpoints: make(map[mappings.Route]mappings.Endpoint),
	}

	for routeName, yamlEndpoint := range serviceConfig.Endpoints {
		route := mappings.Route(routeName)
		endpoint := mappings.Endpoint{
			Method: mappings.Method(yamlEndpoint.Method),
			URI:    yamlEndpoint.URI,
		}
		serviceMapping.Endpoints[route] = endpoint
	}

	return serviceMapping, nil
}

// GetAvailableCountries returns all available countries
func (cl *ConfigLoader) GetAvailableCountries() []countries.Country {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	if cl.config == nil {
		return nil
	}

	countriesList := make([]countries.Country, 0, len(cl.config.Countries))
	for countryStr := range cl.config.Countries {
		countriesList = append(countriesList, countries.Country(countryStr))
	}

	return countriesList
}

// ReloadConfig reloads the configuration from the file
func (cl *ConfigLoader) ReloadConfig() error {
	return cl.LoadConfig()
}

// ValidateConfig validates the loaded configuration
func (cl *ConfigLoader) ValidateConfig() error {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	if cl.config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	// Validate that all countries reference existing services
	for countryName, countryConfig := range cl.config.Countries {
		if _, exists := cl.config.Services[countryConfig.Service]; !exists {
			return fmt.Errorf("country %s references non-existent service %s", countryName, countryConfig.Service)
		}
	}

	// Validate that all services have required endpoints
	requiredEndpoints := []string{"GetAccount"} // Add more as needed
	for serviceName, service := range cl.config.Services {
		for _, requiredEndpoint := range requiredEndpoints {
			if _, exists := service.Endpoints[requiredEndpoint]; !exists {
				return fmt.Errorf("service %s missing required endpoint %s", serviceName, requiredEndpoint)
			}
		}
	}

	return nil
}
