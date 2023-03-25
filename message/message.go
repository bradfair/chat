package message

import (
	"encoding/json"
	"fmt"
	"github.com/bradfair/chat/tokenizer"
)

// Message is a piece of content sent from a role.
type Message struct {
	role      Role
	content   string
	tokenizer tokenizer.Tokenizer
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

// New creates a new message.
func New(options ...Option) Message {
	m := Message{}
	for _, option := range options {
		option(&m)
	}

	return m
}

// Option is a function that configures a message.
type Option func(*Message)

// WithRole configures a message with a role.
func WithRole(role Role) Option {
	return func(m *Message) {
		m.role = role
	}
}

// WithContent configures a message with content.
func WithContent(content string) Option {
	return func(m *Message) {
		m.content = content
	}
}

// WithTokenizer configures a message with a tokenizer.
func WithTokenizer(t tokenizer.Tokenizer) Option {
	return func(m *Message) {
		m.tokenizer = t
	}
}

// ErrNoTokenizer is returned when a message is missing a tokenizer.
var ErrNoTokenizer = fmt.Errorf("missing tokenizer")
