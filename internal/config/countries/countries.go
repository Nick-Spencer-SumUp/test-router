package countries

import (
	"fmt"

	"github.com/Nick-Spencer-SumUp/test-router/internal/config/mappings"
)

type (
	Country       string
	CountryConfig = mappings.ServiceMapping
)

const (
	US Country = "US"
	DE Country = "DE"
)

var (
	CountryMap = map[string]Country{
		"US": US,
		"DE": DE,
	}
)

func GetCountry(country string) (Country, error) {
	result, ok := CountryMap[country]
	if !ok {
		return "", fmt.Errorf("country not found")
	}
	return result, nil
}
