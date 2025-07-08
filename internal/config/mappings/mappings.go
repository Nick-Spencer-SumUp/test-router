package mappings

import "fmt"

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
