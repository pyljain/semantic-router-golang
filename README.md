# Semantic Router

Semantic Router is designed to accelerate decision-making for your services that offer multiple routes and are AI powered. By utilizing semantic vector spaces, it interprets and routes requests based on their meaning, bypassing the delays associated with LLM based tool-use decision making.

This project provides a minimalist semantic routing capability in Golang. It allows you to create multiple routes and map each route to a set of utterances used at runtime for user input classification. These utterances are encoded using an embeddings model and stored in an in-memory index. Currently, the library uses the OpenAI embedding model. For scaling purposes, you might want to fetch these utterances from a configuration store such as a config map in Kubernetes, AWS Config Manager, an object store, or a file persisted in the filesystem. Defining these details is outside the scope of this library, as it focuses on the core logic required for semantic routing.


## Getting started

### Install the library

```sh
go get github.com/pyljain/semantic-router-golang
```


### Usage

An example of how you would define routes and use `semantic_router` to make a decision on which route is the most appropriate based on semantic distance in between the query vector (user input) and these utternances. It returns the closest match computed using a dot product based approach.

```go
routes := []*semantic_router.Route{
    {
        Name: "politics",
        Metadata: semantic_router.RouteMetadata{
            Utterances: []string{
                "isn't politics the best thing ever",
                "why don't you tell me about your political opinions",
                "don't you just love the president",
                "they're going to destroy this country!",
                "they will save the country!",
            },
        },
    },
    {
        Name: "general",
        Metadata: semantic_router.RouteMetadata{
            Utterances: []string{
                "how's the weather today?",
                "how are things going?",
                "lovely weather today",
                "the weather is horrendous",
                "let's go to the chippy",
            },
        },
    },
}

enc, _ := encoder.NewOpenAIEncoder("text-embedding-3-small")

router, _ := semantic_router.New(enc, routes)
result, _ := router.MakeRoutingDecision("How do I make pasta?")
require.Equal(t, "general", result)
/* Include your processing logic once the route determination is made */
```






