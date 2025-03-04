package ollama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

func NewClient(baseUrl string) (*Client, error) {
	var err error

	client := Client{}

	client.baseUrl, err = url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	client.httpClient = http.DefaultClient

	return &client, nil
}

type Client struct {
	baseUrl    *url.URL // baseUrl The base url of the Client server.
	httpClient *http.Client
}

// Chat part

// Message is a single message in a chat sequence. The message contains the
// role ("system", "user", or "assistant"), the content and an optional list
// of images.
type Message struct {
	Role      string      `json:"role"`
	Content   string      `json:"content"`
	Images    []ImageData `json:"images,omitempty"`
	ToolCalls []ToolCall  `json:"tool_calls,omitempty"`
}

// ImageData represents the raw binary data of an image file.
type ImageData []byte

// ChatRequest describes a request sent by [Client.Chat].
type ChatRequest struct {
	// Model is the model name, as in [GenerateRequest].
	Model string `json:"model"`

	// Messages is the messages of the chat - can be used to keep a chat memory.
	Messages []Message `json:"messages"`

	// Stream enables streaming of returned responses; true by default.
	Stream bool `json:"stream"`

	// Format is the format to return the response in (e.g. "json").
	Format json.RawMessage `json:"format,omitempty"`

	// KeepAlive controls how long the model will stay loaded into memory
	// following the request.
	KeepAlive *Duration `json:"keep_alive,omitempty"`

	// Tools is an optional list of tools the model has access to.
	Tools `json:"tools,omitempty"`

	// Options lists model-specific options.
	Options map[string]interface{} `json:"options"`
}

type ChatResponse struct {
	Model      string    `json:"model"`
	CreatedAt  time.Time `json:"created_at"`
	Message    Message   `json:"message"`
	DoneReason string    `json:"done_reason,omitempty"`

	Done bool `json:"done"`

	Metrics
}

type Metrics struct {
	TotalDuration      time.Duration `json:"total_duration,omitempty"`
	LoadDuration       time.Duration `json:"load_duration,omitempty"`
	PromptEvalCount    int           `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration time.Duration `json:"prompt_eval_duration,omitempty"`
	EvalCount          int           `json:"eval_count,omitempty"`
	EvalDuration       time.Duration `json:"eval_duration,omitempty"`
}

func (c *Client) Chat(ctx context.Context, request *ChatRequest) (*ChatResponse, error) {
	var response ChatResponse
	request.Stream = false
	err := c.Post(ctx, "/api/chat", request, &response)
	if err != nil {
		return nil, err
	}

	return &response, err
}

func (c *Client) ChatStream(ctx context.Context, request *ChatRequest, fn func([]byte) error) error {
	request.Stream = true
	err := c.stream(ctx, http.MethodPost, "/api/chat", request, fn)
	if err != nil {
		return err
	}

	return err
}

// Embedding part

// EmbedRequest is the request passed to [Client.Embed].
type EmbedRequest struct {
	// Model is the model name.
	Model string `json:"model"`

	// Input is the input to embed.
	Input any `json:"input"`

	// KeepAlive controls how long the model will stay loaded in memory following
	// this request.
	KeepAlive *Duration `json:"keep_alive,omitempty"`

	Truncate *bool `json:"truncate,omitempty"`

	// Options lists model-specific options.
	Options map[string]interface{} `json:"options"`
}

// EmbedResponse is the response from [Client.Embed].
type EmbedResponse struct {
	Model      string      `json:"model"`
	Embeddings [][]float32 `json:"embeddings"`

	TotalDuration   time.Duration `json:"total_duration,omitempty"`
	LoadDuration    time.Duration `json:"load_duration,omitempty"`
	PromptEvalCount int           `json:"prompt_eval_count,omitempty"`
}

// Embed generates embeddings from a model.
func (c *Client) Embed(ctx context.Context, req *EmbedRequest) (*EmbedResponse, error) {
	var resp EmbedResponse
	err := c.Post(ctx, "/api/embed", req, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// EmbeddingRequest is the request passed to [Client.Embeddings].
type EmbeddingRequest struct {
	// Model is the model name.
	Model string `json:"model"`

	// Prompt is the textual prompt to embed.
	Prompt string `json:"prompt"`

	// KeepAlive controls how long the model will stay loaded in memory following
	// this request.
	KeepAlive *Duration `json:"keep_alive,omitempty"`

	// Options lists model-specific options.
	Options map[string]interface{} `json:"options"`
}

// EmbeddingResponse is the response from [Client.Embeddings].
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// Embeddings generates an embedding from a model.
func (c *Client) Embeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	var resp EmbeddingResponse
	if err := c.Post(ctx, "/api/embeddings", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Request part

func (c *Client) Post(ctx context.Context, path string, data, response any) error {
	httpResp, err := c.do(ctx, http.MethodPost, path, data)
	if err != nil {
		return err
	}

	return c.processResponse(httpResp, response)
}

func (c *Client) do(ctx context.Context, method, path string, data any) (*http.Response, error) {
	var body io.Reader
	var err error
	switch reqData := data.(type) {
	case io.Reader:
		body = reqData
	case nil:
	default:
		var jsonData []byte
		jsonData, err = c.marshalJSON(data)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(jsonData)
	}

	request, err := http.NewRequestWithContext(
		ctx, method, c.baseUrl.JoinPath(path).String(), body,
	)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", fmt.Sprintf("GoAI (%s %s) Go/%s", runtime.GOARCH, runtime.GOOS, runtime.Version()))

	return c.httpClient.Do(request)
}

const maxBufferSize = 512 * 1000

func (c *Client) stream(ctx context.Context, method, path string, data any, fn func([]byte) error) error {
	httpResp, err := c.do(ctx, method, path, data)
	if err != nil {
		return err
	}

	defer httpResp.Body.Close()

	scanner := bufio.NewScanner(httpResp.Body)
	// increase the buffer size to avoid running out of space
	scanBuf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(scanBuf, maxBufferSize)
	for scanner.Scan() {
		var errorResponse struct {
			Error string `json:"error,omitempty"`
		}

		bts := scanner.Bytes()
		if err := json.Unmarshal(bts, &errorResponse); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}

		if errorResponse.Error != "" {
			return errors.New(errorResponse.Error)
		}

		if httpResp.StatusCode >= http.StatusBadRequest {
			return fmt.Errorf("http status code: %s", httpResp.Status)
		}

		if err := fn(bts); err != nil {
			return err
		}
	}

	return nil
}

type errorResponse struct {
	Error string `json:"error"`
}

func (c *Client) processResponse(httpResp *http.Response, response any) error {
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	if httpResp.StatusCode < http.StatusBadRequest {
		return c.unMarshalJSON(respBody, response)
	}

	var e errorResponse
	err = c.unMarshalJSON(respBody, &e)
	if err != nil {
		return fmt.Errorf("http status code: %s, %w", httpResp.Status, err)
	}

	return fmt.Errorf("http status code: %s, %s", httpResp.Status, e.Error)
}

func (c *Client) marshalJSON(data any) ([]byte, error) {
	return json.Marshal(data)
}

func (c *Client) unMarshalJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
