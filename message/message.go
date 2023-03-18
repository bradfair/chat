package message

import "encoding/json"

// Message is a piece of content sent from a role.
type Message struct {
	role    Role
	content string
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
func New(role Role, content string) Message {
	return Message{role: role, content: content}
}
