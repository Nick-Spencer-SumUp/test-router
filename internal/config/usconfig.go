package config

var usConfig = CountryConfig{
	BaseURL: "https://api.atomic.com",
	Endpoints: map[Route]Endpoint{
		GetAccountRoute: {
			Method: GET,
			URI:    "/accounts",
		},
	},
}
