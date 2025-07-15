package config

import (
	"fmt"
	"os"
)

type (
	Service struct {
		Name    string
		BaseUrl string
	}
	
)

var (
	AtomicBaseUrl = getFromEnv("ATOMIC_BASE_URL")
	UpvestBaseUrl = getFromEnv("UPVEST_BASE_URL")

	AtomicService Service = Service{
		Name:    "atomic",
		BaseUrl: AtomicBaseUrl,
	}

	UpvestService Service = Service{
		Name:    "upvest",
		BaseUrl: UpvestBaseUrl,
	}

	CountryServiceMapping = map[string]Service{
		"US": AtomicService,
		"DE": UpvestService,
	}
)

func SelectConfig(country string) (Service, error) {
	service, ok := CountryServiceMapping[country]
	if !ok {
		return Service{}, fmt.Errorf("country %s not supported", country)
	}
	return service, nil
}

func getFromEnv(key string) string {
	return os.Getenv(key)
}
