package tokenizer

import "errors"

type Tokenizer interface {
	// Tokenize splits the given string into a slice of tokens.
	Tokenize(string) ([]int, error)
}

// ErrNoTokenizer is returned when a message has no tokenizer.
var ErrNoTokenizer = errors.New("no tokenizer")

// TokenizerFunc wraps a function as a tokenizer.
type TokenizerFunc func(string) ([]int, error)

// Tokenize calls the wrapped function.
func (f TokenizerFunc) Tokenize(s string) ([]int, error) {
	return f(s)
}
