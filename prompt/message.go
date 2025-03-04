package prompt

const (
	MessageTypeUser      MessageType = "user"
	MessageTypeAssistant MessageType = "assistant"
	MessageTypeSystem    MessageType = "system"
	MessageTypeTool      MessageType = "tool"
)

// MessageType Enumeration representing types of Message in a chat application.
type MessageType string

func (mt MessageType) String() string {
	return string(mt)
}

type Message interface {
	// Type the message type.
	Type() MessageType

	// Text the content of the message
	Text() string

	// Metadata the metadata associated with the content.
	Metadata() map[string]any
}

type defaultMessage struct {
	_type    MessageType
	text     string
	metadata map[string]any
}

func (m *defaultMessage) Type() MessageType {
	return m._type
}

func (m *defaultMessage) Text() string {
	return m.text
}

func (m *defaultMessage) Metadata() map[string]any {
	return m.metadata
}

// UserMessage a message of the type 'user'
func UserMessage(message string) *defaultMessage {
	return &defaultMessage{
		_type: MessageTypeUser,
		text:  message,
	}
}

// AssistantMessage a message of the type 'assistant'
func AssistantMessage(message string) *defaultMessage {
	return &defaultMessage{
		_type: MessageTypeAssistant,
		text:  message,
	}
}

// SystemMessage a message of the type 'system'
func SystemMessage(message string) *defaultMessage {
	return &defaultMessage{
		_type: MessageTypeSystem,
		text:  message,
	}
}

// ToolMessage a message of the type 'tool'
func ToolMessage(message string) *defaultMessage {
	return &defaultMessage{
		_type: MessageTypeTool,
		text:  message,
	}
}
