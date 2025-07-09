package config

import (
	"testing"

	"github.com/Nick-Spencer-SumUp/test-router/internal/config/mappings"
)

func TestYAMLConfigStructures(t *testing.T) {
	// Test YAMLConfig structure
	yamlConfig := YAMLConfig{
		Services: map[string]YAMLService{
			"atomic": {
				BaseURL: "https://api.atomic.com",
				Endpoints: map[string]YAMLEndpoint{
					"GetAccount": {
						Method: "GET",
						URI:    "/accounts",
					},
				},
			},
		},
		Countries: map[string]YAMLCountry{
			"US": {
				Service:  "atomic",
				Region:   "north_america",
				Features: []string{"real_time_payments"},
			},
		},
	}

	// Test that structures are properly defined
	if yamlConfig.Services["atomic"].BaseURL != "https://api.atomic.com" {
		t.Errorf("Expected BaseURL to be https://api.atomic.com, got %s", yamlConfig.Services["atomic"].BaseURL)
	}

	if yamlConfig.Countries["US"].Service != "atomic" {
		t.Errorf("Expected US to use atomic service, got %s", yamlConfig.Countries["US"].Service)
	}

	// Test conversion logic
	loader := NewConfigLoader("")
	loader.config = &yamlConfig

	// Test service mapping conversion
	serviceMapping := mappings.ServiceMapping{
		BaseURL:   yamlConfig.Services["atomic"].BaseURL,
		Endpoints: make(map[mappings.Route]mappings.Endpoint),
	}

	for routeName, yamlEndpoint := range yamlConfig.Services["atomic"].Endpoints {
		route := mappings.Route(routeName)
		endpoint := mappings.Endpoint{
			Method: mappings.Method(yamlEndpoint.Method),
			URI:    yamlEndpoint.URI,
		}
		serviceMapping.Endpoints[route] = endpoint
	}

	// Test endpoint retrieval
	endpointConfig, err := serviceMapping.GetEndpointConfig(mappings.GetAccountRoute)
	if err != nil {
		t.Fatalf("Failed to get endpoint config: %v", err)
	}

	if endpointConfig.Method != mappings.GET {
		t.Errorf("Expected GET method, got %s", endpointConfig.Method)
	}

	if endpointConfig.Endpoint != "/accounts" {
		t.Errorf("Expected /accounts endpoint, got %s", endpointConfig.Endpoint)
	}
}

func TestConfigLoader(t *testing.T) {
	loader := NewConfigLoader("test-config.yaml")

	if loader.configPath != "test-config.yaml" {
		t.Errorf("Expected config path to be test-config.yaml, got %s", loader.configPath)
	}

	// Test that loader is properly initialized
	if loader.config != nil {
		t.Error("Expected config to be nil before loading")
	}
}
