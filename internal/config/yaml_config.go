package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// YAMLConfig represents the combined structure from all configuration files
type YAMLConfig struct {
	Services  map[string]YAMLService `yaml:"services"`
	Countries map[string]YAMLCountry `yaml:"countries"`
}

type YAMLService struct {
	Name         string                  `yaml:"name"`
	BaseURL      string                  `yaml:"base_url"`
	Endpoints    map[string]YAMLEndpoint `yaml:"endpoints"`
	Environments map[string]string       `yaml:"environments"`
}

type YAMLEndpoint struct {
	Method string `yaml:"method"`
	URI    string `yaml:"uri"`
}

type YAMLCountry struct {
	Service      string          `yaml:"service"`
	Environments map[string]bool `yaml:"environments"`
}

// Individual file structures
type CountriesFile struct {
	Countries map[string]YAMLCountry `yaml:"countries"`
}

type ServiceFile struct {
	Service YAMLService `yaml:"service"`
}

// ConfigLoader handles loading and parsing YAML configuration from multiple files
type ConfigLoader struct {
	configDir string
	config    *YAMLConfig
	mu        sync.RWMutex
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader(configDir string) *ConfigLoader {
	return &ConfigLoader{
		configDir: configDir,
	}
}

// LoadConfig loads the configuration from multiple YAML files
func (cl *ConfigLoader) LoadConfig() error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	// Get environment
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "dev"
	}

	// Set default config directory
	configDir := cl.configDir
	if configDir == "" {
		configDir = "internal/config"
	}

	// Load all configuration files
	config := &YAMLConfig{
		Services:  make(map[string]YAMLService),
		Countries: make(map[string]YAMLCountry),
	}

	// Load countries (now includes service mappings)
	if err := cl.loadCountries(configDir, config); err != nil {
		return fmt.Errorf("failed to load countries: %w", err)
	}

	// Load services
	if err := cl.loadServices(configDir, config); err != nil {
		return fmt.Errorf("failed to load services: %w", err)
	}

	// Apply environment-specific overrides
	// Check each service for environment-specific base URLs
	for serviceName, service := range config.Services {
		if service.Environments != nil {
			if envBaseURL, exists := service.Environments[env]; exists && envBaseURL != "" {
				// Override base URL with environment-specific one
				service.BaseURL = envBaseURL
				config.Services[serviceName] = service
			}
		}
	}

	// Filter countries based on environment settings
	filteredCountries := make(map[string]YAMLCountry)
	for countryName, country := range config.Countries {
		// Check if country is enabled for this environment
		if country.Environments != nil {
			if enabled, exists := country.Environments[env]; exists && enabled {
				filteredCountries[countryName] = country
			} else if country.Environments == nil {
				// If no environment config, assume enabled
				filteredCountries[countryName] = country
			}
		} else {
			// If no environment config, assume enabled
			filteredCountries[countryName] = country
		}
	}
	config.Countries = filteredCountries

	cl.config = config
	return nil
}

// loadCountries loads the countries.yaml file (now includes service mappings)
func (cl *ConfigLoader) loadCountries(configDir string, config *YAMLConfig) error {
	countriesPath := filepath.Join(configDir, "countries/countries.yaml")
	data, err := os.ReadFile(countriesPath)
	if err != nil {
		return fmt.Errorf("failed to read countries file %s: %w", countriesPath, err)
	}

	var countriesFile CountriesFile
	if err := yaml.Unmarshal(data, &countriesFile); err != nil {
		return fmt.Errorf("failed to parse countries YAML: %w", err)
	}

	// Load countries with their service mappings
	for countryName, countryConfig := range countriesFile.Countries {
		config.Countries[countryName] = countryConfig
	}

	return nil
}

// loadServices loads all service files from the services directory
func (cl *ConfigLoader) loadServices(configDir string, config *YAMLConfig) error {
	servicesDir := filepath.Join(configDir, "services")
	entries, err := os.ReadDir(servicesDir)
	if err != nil {
		return fmt.Errorf("failed to read services directory %s: %w", servicesDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		servicePath := filepath.Join(servicesDir, entry.Name())
		data, err := os.ReadFile(servicePath)
		if err != nil {
			return fmt.Errorf("failed to read service file %s: %w", servicePath, err)
		}

		var serviceFile ServiceFile
		if err := yaml.Unmarshal(data, &serviceFile); err != nil {
			return fmt.Errorf("failed to parse service YAML %s: %w", servicePath, err)
		}

		// Add service to config
		config.Services[serviceFile.Service.Name] = YAMLService{
			Name:         serviceFile.Service.Name,
			BaseURL:      serviceFile.Service.BaseURL,
			Endpoints:    serviceFile.Service.Endpoints,
			Environments: serviceFile.Service.Environments,
		}
	}

	return nil
}

// GetCountryConfig returns the configuration for a specific country
func (cl *ConfigLoader) GetCountryConfig(country string) (CountryConfig, error) {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	if cl.config == nil {
		return CountryConfig{}, fmt.Errorf("configuration not loaded")
	}

	countryConfig, exists := cl.config.Countries[country]
	if !exists {
		return CountryConfig{}, fmt.Errorf("country %s not found", country)
	}

	serviceConfig, exists := cl.config.Services[countryConfig.Service]
	if !exists {
		return CountryConfig{}, fmt.Errorf("service %s not found for country %s", countryConfig.Service, country)
	}

	// Convert YAML config to ServiceMapping
	serviceMapping := ServiceMapping{
		BaseURL:   serviceConfig.BaseURL,
		Endpoints: make(map[Route]Endpoint),
	}

	for routeName, yamlEndpoint := range serviceConfig.Endpoints {
		route := Route(routeName)
		endpoint := Endpoint{
			Method: Method(yamlEndpoint.Method),
			URI:    yamlEndpoint.URI,
		}
		serviceMapping.Endpoints[route] = endpoint
	}

	return serviceMapping, nil
}

// GetAvailableCountries returns all available countries
func (cl *ConfigLoader) GetAvailableCountries() []Country {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	if cl.config == nil {
		return nil
	}

	countriesList := make([]Country, 0, len(cl.config.Countries))
	for countryStr := range cl.config.Countries {
		countriesList = append(countriesList, Country(countryStr))
	}

	return countriesList
}

// ReloadConfig reloads the configuration from the files
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
	requiredEndpoints := []string{"GetAccount", "UpdateAccount", "DeleteAccount"}
	for serviceName, service := range cl.config.Services {
		for _, requiredEndpoint := range requiredEndpoints {
			if _, exists := service.Endpoints[requiredEndpoint]; !exists {
				return fmt.Errorf("service %s missing required endpoint %s", serviceName, requiredEndpoint)
			}
		}
	}

	return nil
}
