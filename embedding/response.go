package embedding

type Response struct {
	Embeddings []Embedding
}

func (r *Response) List() [][]float32 {
	fs := make([][]float32, len(r.Embeddings))
	for i, embedding := range r.Embeddings {
		fs[i] = embedding.Embedding
	}

	return fs
}
