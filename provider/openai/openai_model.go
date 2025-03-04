package openai

import (
	"context"
	"errors"
	"io"

	"github.com/sashabaranov/go-openai"
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

func (chatModel *ChatModel) Chat(ctx context.Context, message string) (string, error) {
	resp, err := chatModel.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: chatModel.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
	})

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (chatModel *ChatModel) Stream(ctx context.Context, message string, fn func([]byte) error) error {
	stream, err := chatModel.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model: chatModel.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
	})

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
