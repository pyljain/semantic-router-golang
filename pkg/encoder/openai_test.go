package encoder

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenAIEncoder(t *testing.T) {
	tt := []struct {
		description             string
		model                   string
		dataToEncode            []string
		expectError             bool
		expectedDimensions      int
		expectedTotalEmbeddings int
	}{
		{
			description:             "single element to encode",
			model:                   "text-embedding-ada-002",
			dataToEncode:            []string{"How do I get from London to India?"},
			expectError:             false,
			expectedDimensions:      1536,
			expectedTotalEmbeddings: 1,
		},
		{
			description:             "error when no data is passed",
			model:                   "text-embedding-ada-002",
			dataToEncode:            []string{""},
			expectError:             false,
			expectedDimensions:      0,
			expectedTotalEmbeddings: 0,
		},
		{
			description:             "multiple strings passed",
			model:                   "text-embedding-3-large",
			dataToEncode:            []string{"Tell me your name", "my name is Dodgy"},
			expectError:             false,
			expectedDimensions:      3072,
			expectedTotalEmbeddings: 2,
		},
		{
			description:             "total elements are greater than default batch size",
			model:                   "text-embedding-ada-002",
			dataToEncode:            []string{"Tell me your name", "my name is Dodgy", "my name is Wodgy", "my name is Boobit", "my name is Parthurnax", "my name is Minako"},
			expectError:             false,
			expectedDimensions:      1536,
			expectedTotalEmbeddings: 6,
		},
	}

	for _, test := range tt {
		t.Run(test.description, func(t *testing.T) {
			enc, err := NewOpenAIEncoder(test.model)
			if err != nil {
				log.Printf("Error is %s", err)
				if !test.expectError {
					t.Errorf("does not match expected behaviour %s", err)
				}
				return
			}

			resp, err := enc.Encode(test.dataToEncode)
			if err != nil {
				log.Printf("Error is %s", err)
				if !test.expectError {
					t.Errorf("expected no error but received an error %s", err)
				}
				return
			}

			if test.expectedTotalEmbeddings > 0 {
				require.Equal(t, test.expectedDimensions, len(resp[0]))
			}

			require.Equal(t, test.expectedTotalEmbeddings, len(resp))
		})

	}
}
