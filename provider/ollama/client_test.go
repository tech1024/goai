package ollama

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func _handlerFunc(t *testing.T, wantCode int, wantResp any) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(wantCode)
		w.Header().Set("Content-Type", "application/json")
		var err error

		switch wantResp.(type) {
		case string:
			_, err = w.Write([]byte(wantResp.(string)))
		case []byte:
			_, err = w.Write(wantResp.([]byte))
		default:
			err = json.NewEncoder(w).Encode(wantResp)
		}

		if err != nil {
			t.Fatal("failed to encode error response:", err)
		}
	}
}

func TestClient_Chat(t *testing.T) {
	tests := []struct {
		name     string
		request  *ChatRequest
		response *ChatResponse
		wantCode int
		wantErr  error
	}{
		{
			name: "test chat ok",
			request: &ChatRequest{
				Model: "test-model", Messages: []Message{{Role: "user", Content: "request 1"}},
			},
			response: &ChatResponse{
				Model: "test-model", Message: Message{Role: "user", Content: "response 1"},
			},
			wantCode: http.StatusOK,
			wantErr:  nil,
		},
		{
			name:     "test chat model not exist",
			request:  &ChatRequest{},
			response: &ChatResponse{},
			wantCode: http.StatusNotFound,
			wantErr:  errors.New("model not found"),
		},
	}
	var handlerFunc func(writer http.ResponseWriter, request *http.Request)
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handlerFunc(writer, request)
	}))
	defer ts.Close()
	c := &Client{
		baseUrl:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
		httpClient: http.DefaultClient,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc = _handlerFunc(t, tt.wantCode, fmt.Sprintf(`{"model": "%s", "message": {"role": "%s", "content": "%s"}}`,
				tt.response.Model, tt.response.Message.Role, tt.response.Message.Content))
			if tt.wantCode > http.StatusBadRequest {
				handlerFunc = _handlerFunc(t, tt.wantCode, fmt.Sprintf(`{"error": "%s"}`, tt.wantErr.Error()))
			}

			got, err := c.Chat(context.Background(), tt.request)
			if !strings.Contains(fmt.Sprintf("%v", err), fmt.Sprintf("%v", tt.wantErr)) {
				t.Errorf("Chat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tt.response) {
				t.Errorf("Chat() got = %v, want %v", got, tt.response)
			}

			t.Logf("Chat() got = %v", got)
		})
	}
}

func TestClient_Embed(t *testing.T) {
	type args struct {
		ctx context.Context
		req *EmbedRequest
	}
	tests := []struct {
		name     string
		request  *EmbedRequest
		response *EmbedResponse
		wantCode int
		wantErr  error
	}{
		{
			name:     "test embed ok",
			request:  &EmbedRequest{Model: "nomic-embed-text", Input: []string{"hello ai"}},
			response: &EmbedResponse{Model: "nomic-embed-text", Embeddings: [][]float32{{0.25, 0.36}}},
			wantCode: http.StatusOK,
			wantErr:  nil,
		},
		{
			name:     "test embed model not exist",
			request:  &EmbedRequest{},
			response: &EmbedResponse{},
			wantCode: http.StatusNotFound,
			wantErr:  errors.New("model not found"),
		},
	}

	var handlerFunc func(writer http.ResponseWriter, request *http.Request)
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handlerFunc(writer, request)
	}))
	defer ts.Close()
	c := &Client{
		baseUrl:    &url.URL{Scheme: "http", Host: ts.Listener.Addr().String()},
		httpClient: http.DefaultClient,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerFunc = _handlerFunc(t, tt.wantCode, fmt.Sprintf(`{"model": "%s", "embeddings": %s}`,
				tt.response.Model, strings.ReplaceAll(fmt.Sprintf("%v", tt.response.Embeddings), " ", ",")))
			if tt.wantCode > http.StatusBadRequest {
				handlerFunc = _handlerFunc(t, tt.wantCode, fmt.Sprintf(`{"error": "%s"}`, tt.wantErr.Error()))
			}

			got, err := c.Embed(context.Background(), tt.request)
			if !strings.Contains(fmt.Sprintf("%v", err), fmt.Sprintf("%v", tt.wantErr)) {
				t.Errorf("Embed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tt.response) {
				t.Errorf("Embed() got = %v, want %v", got, tt.response)
			}

			t.Logf("Embed() got = %v", got)
		})
	}
}
