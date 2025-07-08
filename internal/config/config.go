package config

import (
	"fmt"
	"github.com/Nick-Spencer-SumUp/test-router/internal/config/countries"
	"github.com/Nick-Spencer-SumUp/test-router/internal/config/mappings"
)

type RoutesConfig map[countries.Country]countries.CountryConfig

var (
	RouterConfigs = RoutesConfig{
		countries.US: countries.USConfig,
		countries.DE: countries.DEConfig,
	}
)

func SelectConfig(country countries.Country, route mappings.Route) (countries.CountryConfig, error) {
	countryConfig, ok := RouterConfigs[country]
	if !ok {
		return countries.CountryConfig{}, fmt.Errorf("country not supported")
	}

	return countryConfig, nil
}
