package index

import (
	"log"
	"math"
	"sort"
)

type InMemoryStore struct {
	data []InsertRequest
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{}
}

func (ims *InMemoryStore) Insert(data []InsertRequest) error {
	ims.data = append(ims.data, data...)
	return nil
}

func (ims *InMemoryStore) Search(query []float64, topK int) ([]SearchResult, error) {
	sm := ims.similarityMatrix(query)
	sort.Slice(sm, func(i, j int) bool {
		return sm[i].Distance > sm[j].Distance
	})

	log.Printf("Sorted sm %+v", sm)

	if len(sm) < topK {
		return sm, nil
	}

	return sm[0:topK], nil
}

func (ims *InMemoryStore) similarityMatrix(query []float64) []SearchResult {
	queryNorm := norm(query)
	result := []SearchResult{}

	for _, ir := range ims.data {
		vec := ir.Embeddings
		vecNorm := norm(vec)

		dotProduct := float64(0)
		for i, e := range vec {
			dotProduct += e * query[i]
		}

		if queryNorm*vecNorm == 0 {
			dotProduct = 0
		} else {
			dotProduct = dotProduct / (queryNorm * vecNorm)
		}

		result = append(result, SearchResult{
			Route:    ir.Route,
			Distance: dotProduct,
		})
	}

	return result
}

func norm(vec []float64) float64 {
	sum := float64(0)
	for _, e := range vec {
		sum += e * e
	}
	return math.Sqrt(sum)
}
