package semantic_router

type Route struct {
	Name     string        `json:"name"`
	Metadata RouteMetadata `json:"metadata"`
}

type RouteMetadata struct {
	Utterances []string `json:"utterances"`
	Category   []string `json:"category"`
}

func NewRoute(name string, utterances []string) *Route {
	return &Route{
		Name: name,
		Metadata: RouteMetadata{
			Utterances: utterances,
		},
	}
}


