package index

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInMemoryStore(t *testing.T) {
	tt := []struct {
		description    string
		query          []float64
		topK           int
		index          []InsertRequest
		expectedResult []SearchResult
	}{
		{
			description:    "When no elements are present in index an empty slice is returned",
			query:          []float64{1, 2, 3},
			topK:           5,
			index:          []InsertRequest{},
			expectedResult: []SearchResult{},
		},
		{
			description: "When elements are present, the result should be sorted (descending) by distance with topK of 1",
			query:       []float64{1, 2, 3},
			topK:        1,
			index: []InsertRequest{
				{
					Route:      "/books",
					Embeddings: []float64{-1, -2, -3},
				},
				{
					Route:      "/cars",
					Embeddings: []float64{1, 2, 3},
				},
				{
					Route:      "/books",
					Embeddings: []float64{0, 0, 0},
				},
				{
					Route:      "/movies",
					Embeddings: []float64{3, 2, 1},
				},
				{
					Route:      "/cars",
					Embeddings: []float64{1, 0, 0},
				},
			},
			expectedResult: []SearchResult{
				{
					Route:    "/cars",
					Distance: 1.0,
				},
			},
		},
		{
			description: "When elements are present, the result should be sorted (descending) by distance with topK of 3",
			query:       []float64{1, 2, 3},
			topK:        3,
			index: []InsertRequest{
				{
					Route:      "/books",
					Embeddings: []float64{-1, -2, -3},
				},
				{
					Route:      "/cars",
					Embeddings: []float64{1, 2, 3},
				},
				{
					Route:      "/books",
					Embeddings: []float64{0, 0, 0},
				},
				{
					Route:      "/movies",
					Embeddings: []float64{3, 2, 1},
				},
				{
					Route:      "/cars",
					Embeddings: []float64{1, 0, 0},
				},
			},
			expectedResult: []SearchResult{
				{
					Route:    "/cars",
					Distance: 1.0,
				},
				{
					Route:    "/movies",
					Distance: 0.7142857142857143,
				},
				{
					Route:    "/cars",
					Distance: 0.2672612419124244,
				},
			},
		},
	}

	for _, test := range tt {
		t.Run(test.description, func(t *testing.T) {
			ims := NewInMemoryStore()
			err := ims.Insert(test.index)
			require.NoError(t, err)

			sr, err := ims.Search(test.query, test.topK)
			require.NoError(t, err)

			for i, res := range sr {
				require.Equal(t, test.expectedResult[i].Distance, res.Distance)
				require.Equal(t, test.expectedResult[i].Route, res.Route)
			}
		})
	}
}
