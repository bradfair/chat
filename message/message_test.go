package message_test

import (
	"github.com/bradfair/chat/message"
	"testing"
)

func TestMessage(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		msg := message.Message{}
		if !msg.IsEmpty() {
			t.Errorf("expected message to be empty")
		}
	})
	t.Run("new", func(t *testing.T) {
		msg := message.New(message.RoleUser, "hello")
		if msg.IsEmpty() {
			t.Errorf("expected message to not be empty")
		}
		if msg.Role() != "user" {
			t.Errorf("expected role to be user")
		}
		if msg.Content() != "hello" {
			t.Errorf("expected content to be hello")
		}
	})
	t.Run("marshal", func(t *testing.T) {
		msg := message.New(message.RoleUser, "hello")
		b, err := msg.MarshalJSON()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if string(b) != `{"role":"user","content":"hello"}` {
			t.Errorf("expected json to be {\"role\":\"user\",\"content\":\"hello\"}, got %s", string(b))
		}
	})
}
