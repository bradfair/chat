package conversation_test

import (
	"errors"
	"github.com/bradfair/chat/conversation"
	"reflect"
	"strings"
	"testing"
)

func TestConversation(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		c := conversation.Conversation{}
		if len(c.Messages()) != 0 {
			t.Errorf("expected conversation to be empty")
		}
	})
	t.Run("new", func(t *testing.T) {
		c := conversation.New()
		if len(c.Messages()) != 0 {
			t.Errorf("expected conversation to be empty")
		}
	})
	t.Run("message out of range", func(t *testing.T) {
		c := conversation.New()
		if c.Message(0) != nil {
			t.Errorf("expected message to be nil")
		}
	})
	t.Run("append", func(t *testing.T) {
		c := conversation.New()
		c.Append(testMessage{role: "user", content: "message 1"})
		c.Append(testMessage{role: "user", content: "message 2"})
		if len(c.Messages()) != 2 {
			t.Errorf("expected conversation to have two messages")
		}
		if c.Message(1).Content() != "message 2" {
			t.Errorf("expected second message to be message 2")
		}
	})
	t.Run("prepend", func(t *testing.T) {
		c := conversation.New()
		c.Prepend(testMessage{role: "user", content: "message 2"})
		c.Prepend(testMessage{role: "user", content: "message 1"})
		if len(c.Messages()) != 2 {
			t.Errorf("expected conversation to have two messages")
		}
		if c.Message(0).Content() != "message 1" {
			t.Errorf("expected first message to be message 1")
		}
	})
	t.Run("remove", func(t *testing.T) {
		c := conversation.New()
		c.Append(testMessage{role: "user", content: "message 1"})
		c.Append(testMessage{role: "user", content: "message 2"})
		m := c.Remove(0)
		if len(c.Messages()) != 1 {
			t.Errorf("expected conversation to have one message")
		}
		if c.Message(0).Content() != "message 2" {
			t.Errorf("expected first message to be message 2")
		}
		if m.Content() != "message 1" {
			t.Errorf("expected removed message to be message 1")
		}
	})
	t.Run("remove out of range", func(t *testing.T) {
		c := conversation.New()
		c.Append(testMessage{role: "user", content: "message 1"})
		c.Append(testMessage{role: "user", content: "message 2"})
		m := c.Remove(2)
		if len(c.Messages()) != 2 {
			t.Errorf("expected conversation to have two messages")
		}
		if m != nil {
			t.Errorf("expected removed message to be nil")
		}
	})
	t.Run("insert", func(t *testing.T) {
		c := conversation.New()
		c.Append(testMessage{role: "user", content: "message 1"})
		c.Append(testMessage{role: "user", content: "message 3"})
		c.Insert(1, testMessage{role: "user", content: "message 2"})
		if len(c.Messages()) != 3 {
			t.Errorf("expected conversation to have three messages")
		}
		if c.Message(1).Content() != "message 2" {
			t.Errorf("expected second message to be message 2")
		}
	})
	t.Run("insert out of range appends", func(t *testing.T) {
		c := conversation.New()
		c.Append(testMessage{role: "user", content: "message 1"})
		c.Append(testMessage{role: "user", content: "message 2"})
		c.Insert(3, testMessage{role: "user", content: "message 3"})
		if len(c.Messages()) != 3 {
			t.Errorf("expected conversation to have three messages")
		}
		if c.Message(2).Content() != "message 3" {
			t.Errorf("expected third message to be message 3")
		}
	})
	t.Run("replace", func(t *testing.T) {
		c := conversation.New()
		c.Append(testMessage{role: "user", content: "message 1"})
		c.Append(testMessage{role: "user", content: "message two"})
		c.Replace(1, testMessage{role: "user", content: "message 2"})
		if len(c.Messages()) != 2 {
			t.Errorf("expected conversation to have two messages")
		}
		if c.Message(1).Content() != "message 2" {
			t.Errorf("expected second message to be message 2")
		}
	})
	t.Run("replace out of range appends", func(t *testing.T) {
		c := conversation.New()
		c.Append(testMessage{role: "user", content: "message 1"})
		c.Append(testMessage{role: "user", content: "message 2"})
		c.Replace(3, testMessage{role: "user", content: "message 3"})
		if len(c.Messages()) != 3 {
			t.Errorf("expected conversation to have three messages")
		}
		if c.Message(2).Content() != "message 3" {
			t.Errorf("expected third message to be message 3")
		}
	})
	t.Run("marshal", func(t *testing.T) {
		c := conversation.New()
		c.Append(testMessage{role: "user", content: "message 1"})
		c.Append(testMessage{role: "user", content: "message 2"})
		b, err := c.MarshalJSON()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if string(b) != `[{"role":"user","content":"message 1"},{"role":"user","content":"message 2"}]` {
			t.Errorf("expected json to be [{\"role\":\"user\",\"content\":\"message 1\"},{\"role\":\"user\",\"content\":\"message 2\"}], got %s", string(b))
		}
	})
	t.Run("parent/child", func(t *testing.T) {
		p := conversation.New()
		p.Append(testMessage{role: "user", content: "message 1"})
		p.Append(testMessage{role: "user", content: "message 2"})
		c := p.NewChild(p.Message(0))
		c.Append(testMessage{role: "user", content: "message 3"})
		c.Append(testMessage{role: "user", content: "message 4"})
		if len(p.Messages()) != 2 {
			t.Errorf("expected parent conversation to have two messages")
		}
		if len(c.Messages()) != 3 {
			t.Errorf("expected child conversation to have three messages")
		}
		if c.Parent() != p {
			t.Errorf("expected parent to be parent")
		}
	})
	t.Run("count tokens", func(t *testing.T) {
		c := conversation.New()
		c.Append(testMessage{role: "user", content: "message 1"})
		c.Append(testMessage{role: "user", content: "message 2"})
		tokenCount, err := c.CountTokens()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if tokenCount != 4 {
			t.Errorf("expected token count to be 4, got %d", tokenCount)
		}
	})
	t.Run("count tokens with error", func(t *testing.T) {
		c := conversation.New()
		c.Append(testMessage{role: "user", content: "message 1"})
		c.Append(testMessage{role: "user", content: "message 2", tokenizeError: errTokenizing})
		tokenCount, err := c.CountTokens()
		if !errors.Is(err, errTokenizing) {
			t.Errorf("expected error to be %v, got %v", errTokenizing, err)
		}
		if tokenCount != 0 {
			t.Errorf("expected token count to be 0, got %d", tokenCount)
		}
	})
	t.Run("WithOptions", func(t *testing.T) {
		t.Run("WithMessages overwrites existing messages", func(t *testing.T) {
			c := conversation.New()
			c.Append(testMessage{role: "user", content: "message 1"})
			want := []conversation.Message{testMessage{role: "user", content: "should overwrite existing messages"}}
			c.WithOptions(conversation.WithMessages(want...))
			if len(c.Messages()) != 1 {
				t.Errorf("expected conversation to have one message")
			}
			if !reflect.DeepEqual(c.Message(0), want[0]) {
				t.Errorf("expected message to be %v, got %v", want[0], c.Message(0))
			}
		})
	})
}

type testMessage struct {
	role          string
	content       string
	tokenizeError error
}

func (m testMessage) Role() string {
	return m.role
}

func (m testMessage) Content() string {
	return m.content
}

func (m testMessage) Tokenize() ([]uint, error) {
	if m.tokenizeError != nil {
		return nil, m.tokenizeError
	}
	var tokens []uint
	for id, _ := range strings.Split(m.content, " ") {
		tokens = append(tokens, uint(id))
	}
	return tokens, nil
}

var errTokenizing = errors.New("error tokenizing")
