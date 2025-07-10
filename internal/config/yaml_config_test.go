package config

import (
	"os"
	"path/filepath"
	"testing"
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
	serviceMapping := ServiceMapping{
		BaseURL:   yamlConfig.Services["atomic"].BaseURL,
		Endpoints: make(map[Route]Endpoint),
	}

	for routeName, yamlEndpoint := range yamlConfig.Services["atomic"].Endpoints {
		route := Route(routeName)
		endpoint := Endpoint{
			Method: Method(yamlEndpoint.Method),
			URI:    yamlEndpoint.URI,
		}
		serviceMapping.Endpoints[route] = endpoint
	}

	// Test endpoint retrieval
	endpointConfig, err := serviceMapping.GetEndpointConfig(GetAccountRoute)
	if err != nil {
		t.Fatalf("Failed to get endpoint config: %v", err)
	}

	if endpointConfig.Method != GET {
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

func TestEnvironmentFiltering(t *testing.T) {
	// Create a temporary directory for test config files
	tempDir := t.TempDir()

	// Create test directory structure
	countriesDir := filepath.Join(tempDir, "countries")
	servicesDir := filepath.Join(tempDir, "services")

	err := os.MkdirAll(countriesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create countries directory: %v", err)
	}

	err = os.MkdirAll(servicesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create services directory: %v", err)
	}

	// Create test countries.yaml with a country only enabled for dev
	countriesYAML := `# Test countries configuration
countries:
  US:
    service: "atomic"
    environments:
      dev: true
      stage: false
      live: false
  DE:
    service: "upvest"
    environments:
      dev: true
      stage: true
      live: true
`

	err = os.WriteFile(filepath.Join(countriesDir, "countries.yaml"), []byte(countriesYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to write test countries.yaml: %v", err)
	}

	// Create test atomic.yaml service file
	atomicYAML := `# Test atomic service configuration
service:
  name: "atomic"
  base_url: "https://api.atomic.com"
  environments:
    dev: "https://dev-api.atomic.com"
    stage: "https://staging-api.atomic.com"
    live: "https://api.atomic.com"
  endpoints:
    GetAccount:
      method: "GET"
      uri: "/accounts"
`

	err = os.WriteFile(filepath.Join(servicesDir, "atomic.yaml"), []byte(atomicYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to write test atomic.yaml: %v", err)
	}

	// Create test upvest.yaml service file
	upvestYAML := `# Test upvest service configuration
service:
  name: "upvest"
  base_url: "https://api.upvest.com"
  environments:
    dev: "https://dev-api.upvest.com"
    stage: "https://staging-api.upvest.com"
    live: "https://api.upvest.com"
  endpoints:
    GetAccount:
      method: "GET"
      uri: "/accounts"
`

	err = os.WriteFile(filepath.Join(servicesDir, "upvest.yaml"), []byte(upvestYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to write test upvest.yaml: %v", err)
	}

	// Test 1: Load config for dev environment - US should be present
	originalEnv := os.Getenv("ENVIRONMENT")
	defer func() {
		os.Setenv("ENVIRONMENT", originalEnv)
	}()

	os.Setenv("ENVIRONMENT", "dev")

	loader := NewConfigLoader(tempDir)
	err = loader.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config for dev environment: %v", err)
	}

	devCountries := loader.GetAvailableCountries()

	// US should be present in dev environment
	var foundUS bool
	for _, country := range devCountries {
		if string(country) == "US" {
			foundUS = true
			break
		}
	}

	if !foundUS {
		t.Error("Expected US to be present in dev environment")
	}

	// Test 2: Load config for stage environment - US should NOT be present
	os.Setenv("ENVIRONMENT", "stage")

	loader2 := NewConfigLoader(tempDir)
	err = loader2.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config for stage environment: %v", err)
	}

	stageCountries := loader2.GetAvailableCountries()

	// US should NOT be present in stage environment
	var foundUSInStage bool
	for _, country := range stageCountries {
		if string(country) == "US" {
			foundUSInStage = true
			break
		}
	}

	if foundUSInStage {
		t.Error("Expected US to NOT be present in stage environment")
	}

	// DE should be present in stage environment (enabled for stage)
	var foundDE bool
	for _, country := range stageCountries {
		if string(country) == "DE" {
			foundDE = true
			break
		}
	}

	if !foundDE {
		t.Error("Expected DE to be present in stage environment")
	}

	// Test 3: Verify we can't get config for filtered country
	_, err = loader2.GetCountryConfig("US")
	if err == nil {
		t.Error("Expected error when trying to get config for filtered country US in stage environment")
	}

	// Test 4: Verify we can get config for non-filtered country
	_, err = loader2.GetCountryConfig("DE")
	if err != nil {
		t.Errorf("Expected to get config for DE in stage environment, got error: %v", err)
	}
}
