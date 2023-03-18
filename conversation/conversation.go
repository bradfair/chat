package conversation

import "encoding/json"

// Conversation is a collection of messages.
type Conversation struct {
	messages []Message
	parent   *Conversation
}

// Messages returns the messages in the conversation.
func (c *Conversation) Messages() []Message {
	c.init()
	return c.messages
}

// Append appends a message to the conversation.
func (c *Conversation) Append(m Message) {
	c.init()
	c.messages = append(c.messages, m)
}

// Prepend prepends a message to the conversation.
func (c *Conversation) Prepend(m Message) {
	c.init()
	c.messages = append([]Message{m}, c.messages...)
}

// Remove removes a message at the given index and returns it. If the index is out of range, nil is returned.
func (c *Conversation) Remove(i uint) Message {
	c.init()
	if i >= uint(len(c.messages)) {
		return nil
	}
	m := c.messages[i]
	c.messages = append(c.messages[:i], c.messages[i+1:]...)
	return m
}

// Insert inserts a message at the given index. If the index is out of range, the message is appended.
func (c *Conversation) Insert(i uint, m Message) {
	c.init()
	if i >= uint(len(c.messages)) {
		c.Append(m)
		return
	}
	c.messages = append(c.messages[:i], append([]Message{m}, c.messages[i:]...)...)
}

// Replace replaces a message at the given index. If the index is out of range, the message is appended.
func (c *Conversation) Replace(i uint, m Message) {
	c.init()
	if i >= uint(len(c.messages)) {
		c.Append(m)
		return
	}
	c.messages[i] = m
}

// MarshalJSON implements the json.Marshaler interface.
func (c *Conversation) MarshalJSON() ([]byte, error) {
	c.init()
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

// Parent returns the parent conversation. If the conversation has no parent, nil is returned.
func (c *Conversation) Parent() *Conversation {
	return c.parent
}

// NewChild returns a new child conversation.
// This is useful for creating a new conversation based on the current one, such as when a chatbot needs to "talk to itself"
// in order to determine its response. Once the chatbot has determined its response, it can easily append the response to
// the parent conversation.
func (c *Conversation) NewChild(messages ...Message) *Conversation {
	return &Conversation{messages: messages, parent: c}
}

func (c *Conversation) init() {
	if c.messages == nil {
		c.messages = []Message{}
	}
}

// New creates a new conversation.
func New(messages ...Message) *Conversation {
	return &Conversation{messages: messages}
}

// Message is a piece of content sent from a role.
type Message interface {
	Role() string
	Content() string
}
