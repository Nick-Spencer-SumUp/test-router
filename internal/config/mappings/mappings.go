package mappings

import (
	"fmt"
	"strings"
)

type (
	Path   string
	Route  = string // Backward compatibility alias
	Method string

	ServiceMapping struct {
		Endpoints map[Path]Endpoint
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
		Path     string
	}
)

const (
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

func (c ServiceMapping) GetEndpointConfigByPath(requestPath string) (EndpointConfig, error) {
	// First try exact match
	if endpoint, ok := c.Endpoints[Path(requestPath)]; ok {
		return EndpointConfig{
			Endpoint: endpoint.URI,
			BaseURL:  c.BaseURL,
			Method:   endpoint.Method,
			Path:     requestPath,
		}, nil
	}

	// If no exact match, try pattern matching for paths with parameters
	for configPath, endpoint := range c.Endpoints {
		if pathMatches(string(configPath), requestPath) {
			return EndpointConfig{
				Endpoint: endpoint.URI,
				BaseURL:  c.BaseURL,
				Method:   endpoint.Method,
				Path:     requestPath,
			}, nil
		}
	}

	return EndpointConfig{}, fmt.Errorf("path %s not supported", requestPath)
}

// pathMatches checks if a request path matches a configured path pattern
// e.g., "/accounts/{id}" matches "/accounts/123"
func pathMatches(configPath, requestPath string) bool {
	configParts := strings.Split(configPath, "/")
	requestParts := strings.Split(requestPath, "/")

	if len(configParts) != len(requestParts) {
		return false
	}

	for i, configPart := range configParts {
		requestPart := requestParts[i]

		// If config part is a parameter (contains {}), it matches anything
		if strings.Contains(configPart, "{") && strings.Contains(configPart, "}") {
			continue
		}

		// Otherwise, parts must match exactly
		if configPart != requestPart {
			return false
		}
	}

	return true
}

// GetAvailablePaths returns all configured paths for this service
func (c ServiceMapping) GetAvailablePaths() []string {
	paths := make([]string, 0, len(c.Endpoints))
	for path := range c.Endpoints {
		paths = append(paths, string(path))
	}
	return paths
}

// Legacy function for backward compatibility - deprecated
func (c ServiceMapping) GetEndpointConfig(route Route) (EndpointConfig, error) {
	return c.GetEndpointConfigByPath(string(route))
}

func (e EndpointConfig) GetMethod() string {
	return MethodMap[e.Method]
}
