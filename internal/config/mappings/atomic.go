package mappings

var AtomicMapping = ServiceMapping{
	BaseURL: "https://api.atomic.com",
	Endpoints: map[Route]Endpoint{
		GetAccountRoute: {
			Method: GET,
			URI:    "/accounts",
		},
	},
}
