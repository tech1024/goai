package openai

import (
	"context"
	"errors"
	"io"

	"github.com/sashabaranov/go-openai"
	"github.com/tech1024/goai/prompt"
)

func NewChatModel(client *openai.Client, model string) *ChatModel {
	return &ChatModel{
		client: client,
		model:  model,
	}
}

type ChatModel struct {
	client *openai.Client
	model  string
}

func (chatModel *ChatModel) Call(ctx context.Context, prompt prompt.Prompt) (string, error) {
	req, err := chatModel.buildChatRequest(prompt)
	if err != nil {
		return "", err
	}

	resp, err := chatModel.client.CreateChatCompletion(ctx, req)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (chatModel *ChatModel) Stream(ctx context.Context, prompt prompt.Prompt, fn func([]byte) error) error {
	req, err := chatModel.buildChatRequest(prompt)
	if err != nil {
		return err
	}

	stream, err := chatModel.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}

	defer stream.Close()

	var resp openai.ChatCompletionStreamResponse
	for {
		resp, err = stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		err = fn([]byte(resp.Choices[0].Delta.Content))
		if err != nil {
			return err
		}
	}

	return nil
}

func (chatModel *ChatModel) buildChatRequest(prompt prompt.Prompt) (openai.ChatCompletionRequest, error) {
	request := openai.ChatCompletionRequest{
		Model:    chatModel.model,
		Messages: make([]openai.ChatCompletionMessage, len(prompt.Messages)),
	}
	for i, message := range prompt.Messages {
		request.Messages[i] = openai.ChatCompletionMessage{
			Role:    message.Type().String(),
			Content: message.Text(),
		}
	}
	if prompt.ChatOption.Model != "" {
		request.Model = prompt.ChatOption.Model
	}

	return request, nil
}
