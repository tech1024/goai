package goai

import (
	"context"

	"github.com/tech1024/goai/prompt"
)

type ChatModel interface {
	Call(ctx context.Context, prompt prompt.Prompt) (string, error)
	Stream(ctx context.Context, prompt prompt.Prompt, receive func([]byte) error) error
}

func NewChat(chatModel ChatModel) *Chat {
	return &Chat{
		chatModel: chatModel,
	}
}

type Chat struct {
	chatModel ChatModel
}

// Chat send a message, it returns string
func (c *Chat) Chat(ctx context.Context, content string) (string, error) {
	return c.Prompt(ctx, prompt.NewPrompt(
		prompt.UserMessage(content),
	))
}

// ChatStream send a message, need to receive its returns
func (c *Chat) ChatStream(ctx context.Context, content string, receive func([]byte) error) error {
	return c.Stream(ctx, prompt.NewPrompt(
		prompt.UserMessage(content),
	), receive)
}

func (c *Chat) Prompt(ctx context.Context, p prompt.Prompt) (string, error) {
	return c.chatModel.Call(ctx, p)
}

func (c *Chat) Stream(ctx context.Context, p prompt.Prompt, fn func([]byte) error) error {
	return c.chatModel.Stream(ctx, p, fn)
}
