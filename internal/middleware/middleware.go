package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GetCountryFromToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, "Authorization required")
		}

		token = strings.TrimPrefix(token, "Bearer ")

		country, err := extractCountryFromToken(token)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid token: %v", err))
		}

		// You can now use the country value or store it in context
		c.Set("country", country)

		return next(c)
	}
}

func extractCountryFromToken(token string) (string, error) {
	parser := new(jwt.Parser)
	parsedToken, _, err := parser.ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to assert claims as MapClaims")
	}

	// Use the helper function to extract the nested value
	merchantCountryValue, err := getNestedValue(claims, "ext", "classic", "merchant_country")
	if err != nil {
		return "", fmt.Errorf("failed to extract merchant_country: %w", err)
	}

	// Final type assertion to string
	merchantCountry, ok := merchantCountryValue.(string)
	if !ok {
		return "", fmt.Errorf("merchant_country is not a string")
	}

	return merchantCountry, nil
}

// Helper function to safely extract nested values from a map
func getNestedValue(m map[string]interface{}, path ...string) (interface{}, error) {
	current := m

	for i, key := range path {
		value, exists := current[key]
		if !exists {
			return nil, fmt.Errorf("key '%s' not found in path %v", key, path[:i+1])
		}

		// If this is the last key in the path, return the value
		if i == len(path)-1 {
			return value, nil
		}

		// Otherwise, ensure the value is a map for further navigation
		nextMap, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("value at '%s' is not a map, cannot navigate further", key)
		}
		current = nextMap
	}

	return nil, fmt.Errorf("empty path provided")
}