package semantic_router

import (
	"log"
	"semantic_router/pkg/encoder"
	"semantic_router/pkg/index"
)

type semanticRouter struct {
	encoder encoder.Encoder
	routes  []*Route
	index   index.Index
}

func New(enc encoder.Encoder, routes []*Route) (*semanticRouter, error) {

	// Instantiate the index
	idx := index.NewInMemoryStore()

	// Encode each utterance and store along with metadata
	for _, r := range routes {
		embeddings, err := enc.Encode(r.Metadata.Utterances)
		if err != nil {
			return nil, err
		}

		reqs := []index.InsertRequest{}
		for _, emd := range embeddings {
			reqs = append(reqs, index.InsertRequest{Route: r.Name, Embeddings: emd})
		}

		err = idx.Insert(reqs)
		if err != nil {
			return nil, err
		}
	}

	return &semanticRouter{
		encoder: enc,
		routes:  routes,
		index:   idx,
	}, nil
}

func (sr *semanticRouter) MakeRoutingDecision(userInput string) (string, error) {

	userInputEmbedding, err := sr.encoder.Encode([]string{userInput})
	if err != nil {
		return "", err
	}

	matches, err := sr.index.Search(userInputEmbedding[0], 3)
	if err != nil {
		return "", err
	}

	log.Printf("matches %v", matches)

	return matches[0].Route, nil
}
