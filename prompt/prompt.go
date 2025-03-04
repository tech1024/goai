package prompt

type Prompt struct {
	Messages   []Message
	ChatOption Option
}

func NewPrompt(messages ...Message) Prompt {
	return Prompt{
		Messages: messages,
	}
}
