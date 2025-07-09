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
				Service: "atomic",
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
	loader := NewConfigLoader("test-config-dir")

	if loader.configDir != "test-config-dir" {
		t.Errorf("Expected config dir to be test-config-dir, got %s", loader.configDir)
	}

	// Test that loader is properly initialized
	if loader.config != nil {
		t.Error("Expected config to be nil before loading")
	}
}

func TestGetCountryFromConfig(t *testing.T) {
	// Test that GetCountryFromConfig returns error when config not initialized
	_, err := GetCountryFromConfig("US")
	if err == nil {
		t.Error("Expected error when config not initialized")
	}
}
