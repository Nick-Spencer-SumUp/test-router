# RFC: Go Microservice Standards for Router/Gateway Services

## Overview

This RFC defines working agreements and standards for Go microservices that function as routers or gateways. These services route requests to downstream services based on country codes.

### Goals
- Establish consistent code standards across services
- Define a scalable project structure
- Ensure maintainable configuration management
- Provide clear patterns for service implementation

### Architecture Principles
- **Configuration-Driven Routing**: Use configuration to determine downstream service routing
	- configs broken into packages (countries, mappings)
		- Mappings contain mappings for our finance services (Atomic/Upvest)
		- Countries allows us to add countries quickly, using 1 of the mappings from above
		- countries.go and mappings.go contain the types and functions relevant to these domains, while upvest.go or de.go contain the actual configured mappings
- **Response Streaming**: Stream responses from downstream services to minimize latency
	- this makes sense to me for efficiency, but also because we don't have the gateways enforcing response types. The normalization of responses would occur in the integration micro-services for this implementation
- **Country-Based Routing**: Route requests based on country codes or regional requirements

### Notes on current config setup
- I have more than 1 place a new country would need added to get a complete config.
	- {country}.go in countries dir
	- typed string in countries.go and the CountryMap
	- country->config registered in RouterConfigs in config.go
- Would like thoughts on how to improve this.
	- considered moving the RouterConfigs and SelectConfig into countries, but they feel a bit like separate domains to me. This would mean we only have to add a new {country}.go in countries though, and register the mapping in countries.go
### Notes on api dir
- **Stand in**: this will be generated instead of hand-rolled

## Current Project Structure

```
service-name/
├── cmd/
│   └── main.go                 # Application entry point
├── api/
│   ├── handlers/              # HTTP handlers (accounts, transfers)
│   │   └── {domain}/
│   │       └── {domain}.go
│   └── routes/                # Route definitions
│       └── {domain}.go
├── internal/
│   ├── config/                # Configuration management
│   │   ├── config.go          # Main configuration logic
│   │   ├── countries/         # Country-specific configurations
│   │   │   ├── countries.go
│   │   │   ├── {country}.go   # e.g., us.go, de.go
│   │   └── mappings/          # Service and endpoint mappings
│   │       ├── mappings.go
│   │       └── {service}.go   # e.g., upvest.go, atomic.go
│   ├── {domain}/              # Business logic (e.g., accounts, transfers)
│   │   ├── service.go         # Service implementation
│   ├── middleware/            # Custom middleware (add country to context)
├── go.mod
├── go.sum
└── README.md
```

## Configuration Management

### Configuration Structure

```go
// config/config.go
type RoutesConfig map[countries.Country]countries.CountryConfig

var (
	RouterConfigs = RoutesConfig{
		countries.US: countries.USConfig,
		countries.DE: countries.DEConfig,
	}
)
```

### Country-Based Configuration

```go
// config/countries/countries.go
type Country string
type CountryConfig = mappings.ServiceMapping

const (
    US Country = "US"
    DE Country = "DE"
    UK Country = "UK"
)

var CountryMap = map[string]Country{
    "US": US,
    "DE": DE,
    "UK": UK,
}
```

### Service Mapping Example
```go
package mappings

var AtomicMapping = ServiceMapping{
	BaseURL: "https://api.atomic.com", // this will come from env variables 
	Endpoints: map[Route]Endpoint{
		GetAccountRoute: {
			Method: GET,
			URI: "/accounts",
		},
	},
}
```
### Service Mapping Pattern

```go
// config/mappings/mappings.go
type ServiceMapping struct {
    Endpoints map[Route]Endpoint `yaml:"endpoints"`
    BaseURL   string             `yaml:"base_url"`
}

type Endpoint struct {
    Method Method `yaml:"method"`
    URI    string `yaml:"uri"`
}
```

### Handlers
	- validate incoming requests, select config, call Service layer using config
	- handle errors with requests, stream responses

```go
// api/handlers/{domain}/{domain}.go
type Handler struct {
    service Service
    logger  *log.Logger
}

func New(service Service, logger *log.Logger) *Handler {
    return &Handler{
        service: service,
        logger:  logger,
    }
}

func (h *Handler) GetResource(c echo.Context) error {
    // Extract country from request
    country := extractCountryFromRequest(c)
    
    // Get configuration
    cfg, err := config.SelectConfig(country, mappings.GetResourceRoute)
    if err != nil {
        return handleError(c, err)
    }
    
    // Parse request
    var req Request
    if err := c.Bind(&req); err != nil {
        return handleError(c, err)
    }
    
    // Call service
    response, err := h.service.GetResource(c.Request().Context(), cfg, req)
    if err != nil {
        return handleError(c, err)
    }
    
    // Stream response
    return streamResponse(c, response)
}
```

### Response Streaming Pattern

```go
func streamResponse(c echo.Context, response *http.Response) error {
    defer response.Body.Close()
    
    // Copy headers
    for key, values := range response.Header {
        for _, value := range values {
            c.Response().Header().Add(key, value)
        }
    }
    
    // Set status code
    c.Response().WriteHeader(response.StatusCode)
    
    // Stream body
    _, err := io.Copy(c.Response().Writer, response.Body)
    return err
}
```
