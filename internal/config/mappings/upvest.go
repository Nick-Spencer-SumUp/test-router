package mappings

var UpvestMapping = ServiceMapping{
	BaseURL: "https://api.upvest.com",
	Endpoints: map[Route]Endpoint{
		GetAccountRoute: {
			Method: GET,
			URI:    "/accounts",
		},
	},
}