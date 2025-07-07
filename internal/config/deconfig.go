package config

var deConfig = CountryConfig{
	BaseURL: "https://api.upvest.com",
	Endpoints: map[Route]Endpoint{
		GetAccountRoute: {
			Method: GET,
			URI:    "/accounts",
		},
	},
}