package countries

import "github.com/Nick-Spencer-SumUp/test-router/internal/config/mappings"

type (
	Country string
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

func GetCountry(country string) Country {
	return CountryMap[country]
}