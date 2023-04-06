package conversation

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// Conversation is a collection of messages.
type Conversation struct {
	messages []Message
	parent   *Conversation
	mutex    sync.Mutex
}

// Messages returns the messages in the conversation.
func (c *Conversation) Messages() Messages {
	c.mutex.Lock()
	c.mutex.Unlock()
	c.init()
	return c.messages
}

// Message returns the message at the given index. If the index is out of range, nil is returned.
func (c *Conversation) Message(i int) Message {
	if i < 0 {
		return nil
	}
	c.mutex.Lock()
	c.mutex.Unlock()
	c.init()
	if i >= len(c.messages) {
		return nil
	}
	return c.messages[i]
}

// CountTokens returns the number of tokens in the conversation.
func (c *Conversation) CountTokens() (int, error) {
	c.mutex.Lock()
	c.mutex.Unlock()
	c.init()
	var count int
	for _, m := range c.messages {
		tokens, err := m.Tokenize()
		if err != nil {
			return 0, fmt.Errorf("could not tokenize message %q: %w", m.Content(), err)
		}
		count += len(tokens)
	}
	return count, nil
}

// Append appends a message to the conversation.
func (c *Conversation) Append(m Message) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.init()
	c.messages = append(c.messages, m)
}

// Prepend prepends a message to the conversation.
func (c *Conversation) Prepend(m Message) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.init()
	c.messages = append([]Message{m}, c.messages...)
}

// Remove removes a message at the given index and returns it. If the index is out of range, nil is returned.
func (c *Conversation) Remove(i int) Message {
	if i < 0 {
		return nil
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.init()
	if i >= len(c.messages) {
		return nil
	}
	m := c.messages[i]
	c.messages = append(c.messages[:i], c.messages[i+1:]...)
	return m
}

// Insert inserts a message at the given index. If the index is out of range, the message is appended.
func (c *Conversation) Insert(i int, m Message) {
	c.mutex.Lock()
	c.init()
	if i >= len(c.messages) || i < 0 {
		c.Append(m)
		c.mutex.Unlock()
		return
	}
	defer c.mutex.Unlock()
	c.messages = append(c.messages[:i], append([]Message{m}, c.messages[i:]...)...)
}

// Replace replaces a message at the given index. If the index is out of range, the message is appended.
func (c *Conversation) Replace(i int, m Message) {
	c.mutex.Lock()
	c.init()
	if i >= len(c.messages) || i < 0 {
		c.Append(m)
		c.mutex.Unlock()
		return
	}
	defer c.mutex.Unlock()
	c.messages[i] = m
}

// MarshalJSON implements the json.Marshaler interface.
func (c *Conversation) MarshalJSON() ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
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
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.parent
}

// NewChild returns a new child conversation.
// This is useful for creating a new conversation based on the current one, such as when a chatbot needs to "talk to itself"
// in order to determine its response. Once the chatbot has determined its response, it can easily append the response to
// the parent conversation.
func (c *Conversation) NewChild() *Conversation {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return New().WithParent(c)
}

// WithMessages sets the messages in the conversation.
func (c *Conversation) WithMessages(messages ...Message) *Conversation {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.init()
	c.messages = messages
	return c
}

// WithParent sets the parent conversation.
func (c *Conversation) WithParent(parent *Conversation) *Conversation {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.parent = parent
	return c
}

// init initializes the conversation.
func (c *Conversation) init() {
	if c.messages == nil {
		c.messages = []Message{}
	}
}

// New creates a new conversation.
func New() *Conversation {
	c := new(Conversation)
	return c
}

// Message is a piece of content sent from a role.
type Message interface {
	Role() string
	Content() string
	Tokenize() ([]int, error)
}

type Messages []Message

// Len returns the number of messages.
func (m Messages) Len() int {
	return len(m)
}

// Transcript returns the messages as a string.
func (m Messages) Transcript() string {
	var transcript string
	for _, message := range m {
		transcript += fmt.Sprintf("%s: %s\n", message.Role(), message.Content())
	}
	return strings.TrimSpace(transcript)
}
