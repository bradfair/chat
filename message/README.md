# Message Package
This package provides a Message struct with relevant methods/functions to create and manipulate messages.

## Usage
### Creating a Message
To create a new message, use the New function along with the relevant Options to configure the message's role, content, and tokenizer:

```go
import "github.com/bradfair/chat/message"

m := message.New(
    message.WithRole("user"),
    message.WithContent("Hello, world!"),
)
```

### Accessing Message Properties
You can access a message's role and content using the Role and Content methods:

```go
role := m.Role()
content := m.Content()
```
### Checking If a Message Is Empty
To check if a message is empty, use the IsEmpty method:

```go
isEmpty := m.IsEmpty()
```

### Tokenizing a Message
To tokenize a message, use the Tokenize method. You will need to configure the message with a tokenizer before calling this method, otherwise `ErrNoTokenizer` will be returned.

```go
m := message.New(
    message.WithRole("user"),
    message.WithContent("Hello, world!"),
    message.WithTokenizer(tokenizerInstance),
)

tokens, err := m.Tokenize()
if err != nil {
    // Handle error
}
```

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
    "github.com/bradfair/chat/message"
    encoder "github.com/samber/go-gpt-3-encoder"
)

enc, err := encoder.NewEncoder()
if err != nil {
    log.Fatal(err)
}

myTokenizer := message.TokenizerFunc(enc.Encode)
```

## License
This package is released under the MIT License. See [LICENSE](/LICENSE) for more information.