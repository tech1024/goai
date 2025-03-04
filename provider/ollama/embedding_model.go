package ollama

import (
	"context"

	"github.com/tech1024/goai/embedding"
)

func NewEmbeddingModel(client *Client, model string) *EmbeddingModel {
	return &EmbeddingModel{
		client: client,
		model:  model,
	}
}

type EmbeddingModel struct {
	client *Client
	model  string
}

func (embeddingModel *EmbeddingModel) Call(ctx context.Context, request embedding.Request) (embedding.Response, error) {
	er := EmbedRequest{
		Model: embeddingModel.model,
		Input: request.Inputs,
	}

	var embeddingResponse embedding.Response
	resp, err := embeddingModel.client.Embed(ctx, &er)

	if err != nil {
		return embeddingResponse, err
	}

	embeddingResponse.Embeddings = make([]embedding.Embedding, len(resp.Embeddings))
	for i, es := range resp.Embeddings {
		embeddingResponse.Embeddings[i] = embedding.Embedding{
			Embedding: es,
			Index:     i,
		}
	}

	return embeddingResponse, nil
}
