package message

import (
	"encoding/json"
)

// Message is a piece of content sent from a role.
type Message struct {
	role      Role
	content   string
	tokenizer Tokenizer
}

// Role returns the name of the role that sent the message.
func (m Message) Role() string {
	return string(m.role)
}

// Content returns the content of the message.
func (m Message) Content() string {
	return m.content
}

// IsEmpty returns true if the message is empty.
func (m Message) IsEmpty() bool {
	return m.role == "" && m.content == ""
}

// Tokenize returns the message content as a slice of tokens.
func (m Message) Tokenize() ([]int, error) {
	if m.tokenizer == nil {
		return nil, ErrNoTokenizer
	}

	return m.tokenizer.Tokenize(m.content)
}

// MarshalJSON implements the json.Marshaler interface.
func (m Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{
		Role:    m.Role(),
		Content: m.Content(),
	})
}

// WithRole configures a message with a role.
func (m Message) WithRole(role Role) Message {
	m.role = role
	return m
}

// WithContent configures a message with content.
func (m Message) WithContent(content string) Message {
	m.content = content
	return m
}

// WithTokenizer configures a message with a tokenizer.
func (m Message) WithTokenizer(t Tokenizer) Message {
	m.tokenizer = t
	return m
}

// New creates a new message.
func New() Message {
	m := Message{}
	return m
}
