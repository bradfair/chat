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

## License
This package is released under the MIT License. See [LICENSE](/LICENSE) for more information.