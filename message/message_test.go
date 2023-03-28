package message_test

import (
	"errors"
	"github.com/bradfair/chat/message"
	"strings"
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
		msg := message.New(message.WithRole("user"), message.WithContent("hello"))
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
		msg := message.New(message.WithRole("user"), message.WithContent("hello"))
		b, err := msg.MarshalJSON()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if string(b) != `{"role":"user","content":"hello"}` {
			t.Errorf("expected json to be {\"role\":\"user\",\"content\":\"hello\"}, got %s", string(b))
		}
	})
	t.Run("tokenize", func(t *testing.T) {
		t.Run("no tokenizer", func(t *testing.T) {
			msg := message.New(message.WithRole("user"), message.WithContent("hello"))
			tokens, err := msg.Tokenize()
			if !errors.Is(err, message.ErrNoTokenizer) {
				t.Errorf("expected ErrNoTokenizer, got %v", err)
			}
			if tokens != nil {
				t.Errorf("expected tokens to be nil")
			}
		})
		t.Run("with tokenizer", func(t *testing.T) {
			msg := message.New(
				message.WithRole("user"),
				message.WithContent("hello world"),
				message.WithTokenizer(message.TokenizerFunc(testTokenizer)),
			)
			tokens, err := msg.Tokenize()
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if len(tokens) != 2 {
				t.Errorf("expected 2 tokens, got %d", len(tokens))
			}
			if tokens[0] != 0 {
				t.Errorf("expected first token to be 0, got %d", tokens[0])
			}
			if tokens[1] != 1 {
				t.Errorf("expected second token to be 1, got %d", tokens[1])
			}
		})
	})
}

func testTokenizer(content string) (tokens []int, err error) {
	for id, _ := range strings.Split(content, " ") {
		tokens = append(tokens, id)
	}
	return
}
