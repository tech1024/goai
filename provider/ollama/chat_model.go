package ollama

import (
	"context"

	"github.com/tech1024/goai/prompt"
)

func NewNewChatModel(client *Client, model string) *ChatModel {
	return &ChatModel{
		client: client,
		model:  model,
	}
}

type ChatModel struct {
	client *Client
	model  string
}

func (chatModel *ChatModel) Call(ctx context.Context, prompt prompt.Prompt) (string, error) {
	req, err := chatModel.buildChatRequest(prompt)
	if err != nil {
		return "", err
	}

	resp, err := chatModel.client.Chat(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Message.Content, nil
}

func (chatModel *ChatModel) Stream(ctx context.Context, prompt prompt.Prompt, fn func([]byte) error) error {
	req, err := chatModel.buildChatRequest(prompt)
	if err != nil {
		return err
	}

	err = chatModel.client.ChatStream(ctx, req, fn)
	if err != nil {
		return err
	}

	return nil
}

func (chatModel *ChatModel) buildChatRequest(prompt prompt.Prompt) (*ChatRequest, error) {
	request := ChatRequest{
		Model:    chatModel.model,
		Messages: make([]Message, len(prompt.Messages)),
	}
	for i, message := range prompt.Messages {
		request.Messages[i] = Message{
			Role:    message.Type().String(),
			Content: message.Text(),
		}
	}
	if prompt.ChatOption.Model != "" {
		request.Model = prompt.ChatOption.Model
	}

	return &request, nil
}
