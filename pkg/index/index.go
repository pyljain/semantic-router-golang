package index

type Index interface {
	Search(query []float64, topK int) ([]SearchResult, error)
	Insert([]InsertRequest) error
}

type SearchResult struct {
	Route    string
	Distance float64
}

type InsertRequest struct {
	Route      string
	Embeddings []float64
}
