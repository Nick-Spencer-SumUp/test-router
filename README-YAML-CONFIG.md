# YAML Configuration System for Router Service

## Overview

This document explains the benefits of using YAML configuration over Go code for ServiceMappings in the router service and provides implementation details.

## Benefits of YAML Configuration

### 1. **Runtime Configuration Changes**
- **Before**: Required recompilation and redeployment for any endpoint changes
- **After**: Update YAML file and call `/admin/reload-config` endpoint

### 2. **Environment-Specific Configuration**
- **Before**: Same hardcoded URLs across all environments
- **After**: Different configs per environment (dev/staging/prod)

### 3. **DevOps/SRE Friendly**
- **Before**: Developers needed to modify Go code for config changes
- **After**: Operations teams can manage configurations independently

### 4. **Centralized Configuration Management**
- **Before**: Configuration scattered across multiple Go files
- **After**: Single YAML file with all routing configuration

### 5. **External Configuration Management**
- **Before**: Configuration baked into binary
- **After**: Can be managed via ConfigMaps, Consul, etc.

## Configuration Structure

### Main Configuration File (`config/routing.yaml`)

```yaml
# Service Mappings Configuration
services:
  atomic:
    base_url: "https://api.atomic.com"
    endpoints:
      GetAccount:
        method: "GET"
        uri: "/accounts"
      CreateAccount:
        method: "POST"
        uri: "/accounts"
      UpdateAccount:
        method: "PUT"
        uri: "/accounts/{id}"
      DeleteAccount:
        method: "DELETE"
        uri: "/accounts/{id}"

  upvest:
    base_url: "https://api.upvest.com"
    endpoints:
      GetAccount:
        method: "GET"
        uri: "/accounts"
      CreateAccount:
        method: "POST"
        uri: "/accounts"

# Country to Service Mapping
countries:
  US:
    service: "atomic"
    region: "north_america"
    features:
      - "real_time_payments"
      - "instant_transfers"
  
  DE:
    service: "upvest"
    region: "europe"
    features:
      - "sepa_transfers"
      - "euro_payments"

# Environment-specific overrides
environments:
  development:
    services:
      atomic:
        base_url: "https://dev-api.atomic.com"
      upvest:
        base_url: "https://dev-api.upvest.com"
  
  staging:
    services:
      atomic:
        base_url: "https://staging-api.atomic.com"
      upvest:
        base_url: "https://staging-api.upvest.com"
  
  production:
    services:
      atomic:
        base_url: "https://api.atomic.com"
      upvest:
        base_url: "https://api.upvest.com"
```

## Usage

### 1. Starting the Service

```bash
# Use default config file (config/routing.yaml)
go run cmd/main.go

# Use custom config file
CONFIG_PATH=/path/to/config.yaml go run cmd/main.go

# Set environment
ENVIRONMENT=staging go run cmd/main.go
```

### 2. Admin Endpoints

The service now includes admin endpoints for configuration management:

```bash
# Reload configuration without restart
curl -X POST http://localhost:8080/admin/reload-config

# Get available countries
curl http://localhost:8080/admin/countries

# Health check
curl http://localhost:8080/admin/health
```

### 3. Environment Variables

- `CONFIG_PATH`: Path to the YAML configuration file (default: `config/routing.yaml`)
- `ENVIRONMENT`: Environment name for configuration overrides (default: `development`)

## Migration from Go Code

### Before (Go Code)
```go
// internal/config/mappings/upvest.go
var UpvestMapping = ServiceMapping{
    BaseURL: "https://api.upvest.com",
    Endpoints: map[Route]Endpoint{
        GetAccountRoute: {
            Method: GET,
            URI:    "/accounts",
        },
    },
}

// internal/config/countries/de.go
var DEConfig = mappings.UpvestMapping

// internal/config/config.go
var RouterConfigs = RoutesConfig{
    countries.US: countries.USConfig,
    countries.DE: countries.DEConfig,
}
```

### After (YAML Configuration)
```yaml
# config/routing.yaml
services:
  upvest:
    base_url: "https://api.upvest.com"
    endpoints:
      GetAccount:
        method: "GET"
        uri: "/accounts"

countries:
  DE:
    service: "upvest"
    region: "europe"
```

## Adding New Countries/Services

### Adding a New Country

1. Add the country to the `countries` section in `routing.yaml`:
```yaml
countries:
  FR:
    service: "upvest"  # Use existing service
    region: "europe"
    features:
      - "sepa_transfers"
```

2. Reload configuration:
```bash
curl -X POST http://localhost:8080/admin/reload-config
```

### Adding a New Service

1. Add the service to the `services` section:
```yaml
services:
  new_service:
    base_url: "https://api.newservice.com"
    endpoints:
      GetAccount:
        method: "GET"
        uri: "/accounts"
```

2. Update country configurations to use the new service:
```yaml
countries:
  UK:
    service: "new_service"
    region: "europe"
```

## Validation

The system includes automatic validation:
- All countries must reference existing services
- All services must have required endpoints
- Configuration is validated on startup and reload

## Benefits Summary

| Aspect | Go Code | YAML Configuration |
|--------|---------|-------------------|
| **Runtime Changes** | ❌ Requires recompilation | ✅ Hot-reload via API |
| **Environment Support** | ❌ Hardcoded values | ✅ Environment-specific overrides |
| **DevOps Friendly** | ❌ Requires code changes | ✅ Configuration-only changes |
| **External Management** | ❌ Baked into binary | ✅ ConfigMaps, Consul, etc. |
| **Centralized Config** | ❌ Scattered across files | ✅ Single YAML file |
| **Type Safety** | ✅ Compile-time validation | ❌ Runtime validation |
| **Performance** | ✅ No parsing overhead | ❌ Minimal parsing overhead |

## Best Practices

1. **Validation**: Always validate configuration after changes
2. **Version Control**: Keep configuration files in version control
3. **Environment Separation**: Use environment-specific overrides
4. **Monitoring**: Monitor configuration reload events
5. **Rollback**: Have a rollback strategy for configuration changes

## Implementation Details

The YAML configuration system:
- Uses `gopkg.in/yaml.v3` for parsing
- Provides thread-safe configuration loading
- Maintains backward compatibility with existing Go types
- Includes comprehensive validation
- Supports environment-specific overrides
- Provides admin endpoints for management

This approach gives you the flexibility of external configuration while maintaining the strong typing and structure of your existing Go code. 