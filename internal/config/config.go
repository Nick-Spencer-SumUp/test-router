package config

import (
	"fmt"
	"os"
)

type (
	Country string
	Service struct {
		Name    string
		BaseUrl string
	}

	CountryServiceMapping map[Country]Service
)

var (
	CountryServiceMap = CountryServiceMapping{}
	US Country = "US"
	DE Country = "DE"
	CountryMap        = map[string]Country{
		"US": US,
		"DE": DE,
	}
)

func (c CountryServiceMapping) SelectConfig(country Country) (Service, error) {
	service, ok := c[country]
	if !ok {
		return Service{}, fmt.Errorf("country %s not supported", country)
	}
	return service, nil
}

func getFromEnv(key string) string {
	return os.Getenv(key)
}

func InitCountryServiceMapping() CountryServiceMapping {
	AtomicBaseUrl := getFromEnv("ATOMIC_BASE_URL")
	UpvestBaseUrl := getFromEnv("UPVEST_BASE_URL")

	AtomicService := Service{
		Name:    "atomic",
		BaseUrl: AtomicBaseUrl,
	}
	UpvestService := Service{
		Name:    "upvest",
		BaseUrl: UpvestBaseUrl,
	}

	CountryServiceMap = CountryServiceMapping{
		US: AtomicService,
		DE: UpvestService,
	}

	return CountryServiceMap
}

func CountryFromString(country string) (Country, error) {
	countryType, ok := CountryMap[country]
	if !ok {
		return Country(""), fmt.Errorf("country %s not supported", country)
	}
	return countryType, nil
}
