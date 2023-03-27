# Tokenizer Package
This package defines a Tokenizer interface and provides a TokenizerFunc type in order to satisfy the Tokenizer interface using custom functions. This allows you to use your own tokenizer function, or one of the small number of available open-source tokenizers.

## Usage
### Implementing a Custom Tokenizer
To create a custom tokenizer, implement the Tokenizer interface:

```go
type MyTokenizer struct {}

func (t MyTokenizer) Tokenize(s string) ([]int, error) {
    // Tokenization logic goes here
}
```

### Using TokenizerFunc
`TokenizerFunc` is a wrapper that allows you to use a function as a tokenizer. This is useful when you want to use an existing function or method as a tokenizer without creating a new struct.

For example, using the Encode method from the [github.com/samber/go-gpt-3-encoder](https://github.com/samber/go-gpt-3-encoder) package:

```go
import (
    "github.com/bradfair/chat/tokenizer"
    encoder "github.com/samber/go-gpt-3-encoder"
)

enc, err := encoder.NewEncoder()
if err != nil {
    log.Fatal(err)
}

myTokenizer := tokenizer.TokenizerFunc(enc.Encode)
```

## License
This package is released under the MIT License. See [LICENSE](/LICENSE) for more information.