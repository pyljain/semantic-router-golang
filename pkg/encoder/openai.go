package encoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	embeddingBaseURL = "https://api.openai.com/v1"
	batchSize        = 5
)

type openAIEncoder struct {
	apiKey string
	model  string
}

func NewOpenAIEncoder(model string) (*openAIEncoder, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("an OpenAI API key is requirement, environment variable OPENAI_API_KEY is not set")
	}
	return &openAIEncoder{
		apiKey: key,
		model:  model,
	}, nil
}

func (oai *openAIEncoder) Encode(data []string) ([][]float64, error) {
	allEmbeddingsForUtterances := [][]float64{}
	batchForEmbedding := []string{}
	for _, ut := range data {
		batchForEmbedding = append(batchForEmbedding, ut)
		if len(batchForEmbedding) == batchSize {
			generatedEmbeddings, err := oai.generateEmbedding(batchForEmbedding)
			if err != nil {
				return nil, err
			}

			allEmbeddingsForUtterances = append(allEmbeddingsForUtterances, generatedEmbeddings...)
			batchForEmbedding = []string{}
		}
	}

	if len(batchForEmbedding) > 0 {
		generatedEmbeddings, err := oai.generateEmbedding(batchForEmbedding)
		if err != nil {
			return nil, err
		}

		allEmbeddingsForUtterances = append(allEmbeddingsForUtterances, generatedEmbeddings...)
	}

	return allEmbeddingsForUtterances, nil
}

func (oai *openAIEncoder) generateEmbedding(batch []string) ([][]float64, error) {
	reqBody := openAIEmbeddingRequest{
		Input:          batch,
		Model:          oai.model,
		EncodingFormat: "float",
	}

	reqB, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(reqB)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/embeddings", embeddingBaseURL), buff)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", oai.apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := openAIEmbeddingResponse{}
	err = json.Unmarshal(respB, &res)
	if err != nil {
		return nil, err
	}

	embeddings := [][]float64{}
	for _, object := range res.Data {
		embeddings = append(embeddings, object.Embedding)
	}

	return embeddings, nil

}

type openAIEmbeddingRequest struct {
	Input          []string `json:"input"`
	Model          string   `json:"model"`
	EncodingFormat string   `json:"encoding_format"`
}

type openAIEmbeddingResponse struct {
	Data []embeddingData `json:"data"`
}

type embeddingData struct {
	Embedding []float64 `json:"embedding"`
}
