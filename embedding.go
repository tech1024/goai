package goai

import (
	"context"

	"github.com/tech1024/goai/embedding"
)

type EmbeddingModel interface {
	Call(context.Context, embedding.Request) (embedding.Response, error)
}

type Embedding struct {
	embeddingModel EmbeddingModel
}

func (e *Embedding) Embed(ctx context.Context, text string) ([]float32, error) {
	es, err := e.Embeds(ctx, text)
	if err != nil {
		return nil, err
	}

	return es[0], nil
}

func (e *Embedding) Embeds(ctx context.Context, texts ...string) ([][]float32, error) {
	response, err := e.embeddingModel.Call(
		ctx,
		embedding.NewRequest(texts, embedding.Option{}),
	)

	if err != nil {
		return nil, err
	}

	return response.List(), nil
}
