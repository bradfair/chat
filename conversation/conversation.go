package conversation

import "encoding/json"

// Conversation is a collection of messages.
type Conversation struct {
	messages []Message
}

// Messages returns the messages in the conversation.
func (c *Conversation) Messages() []Message {
	return c.messages
}

// Append appends a message to the conversation.
func (c *Conversation) Append(m Message) {
	c.messages = append(c.messages, m)
}

// Prepend prepends a message to the conversation.
func (c *Conversation) Prepend(m Message) {
	c.messages = append([]Message{m}, c.messages...)
}

// Remove removes a message at the given index and returns it. If the index is out of range, nil is returned.
func (c *Conversation) Remove(i uint) Message {
	if i >= uint(len(c.messages)) {
		return nil
	}
	m := c.messages[i]
	c.messages = append(c.messages[:i], c.messages[i+1:]...)
	return m
}

// Insert inserts a message at the given index. If the index is out of range, the message is appended.
func (c *Conversation) Insert(i uint, m Message) {
	if i >= uint(len(c.messages)) {
		c.Append(m)
		return
	}
	c.messages = append(c.messages[:i], append([]Message{m}, c.messages[i:]...)...)
}

// Replace replaces a message at the given index. If the index is out of range, the message is appended.
func (c *Conversation) Replace(i uint, m Message) {
	if i >= uint(len(c.messages)) {
		c.Append(m)
		return
	}
	c.messages[i] = m
}

// MarshalJSON implements the json.Marshaler interface.
func (c *Conversation) MarshalJSON() ([]byte, error) {
	type message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	var messages []message
	for _, m := range c.messages {
		messages = append(messages, message{
			Role:    m.Role(),
			Content: m.Content(),
		})
	}
	return json.Marshal(messages)
}

// New creates a new conversation.
func New() *Conversation {
	return &Conversation{}
}

// Message is a piece of content sent from a role.
type Message interface {
	Role() string
	Content() string
}
