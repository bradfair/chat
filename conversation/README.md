# Conversation Package
This package provides a Conversation struct with relevant methods/functions to create and manipulate conversations.

## Usage
### Creating a Conversation
To create a new conversation, use the New function:

```go
import "github.com/bradfair/chat/conversation"

c := conversation.New()
```

### Adding Messages to a Conversation
To add messages to a conversation, use the Append and Prepend methods:

```go
import (
    "github.com/bradfair/chat/conversation"
    "github.com/bradfair/chat/message"
)

c := conversation.New()

m1 := message.New(message.WithRole("user"), message.WithContent("Hello!"))
m2 := message.New(message.WithRole("assistant"), message.WithContent("Hi there!"))

c.Append(m1) // Append a message to the end of the conversation
c.Prepend(m2) // Prepend a message to the beginning of the conversation
```

### Accessing and Modifying Messages in a Conversation
You can access messages, count tokens, insert, remove, or replace messages using the provided methods:

```go
// Get all messages in the conversation
msgs := c.Messages()

// Get a specific message by index
msg := c.Message(1)

// Count tokens in the conversation. This will return an error if any of the messages in the conversation do not have a tokenizer configured.
tokenCount, err := c.CountTokens()

// Remove a message at a specific index. This is useful for removing earlier messages in a conversation in order to reclaim tokens.
removedMsg := c.Remove(1)

// Insert a message at a specific index
c.Insert(1, m1)

// Replace a message at a specific index. This can be used to replace a message with a more useful prompt-based message. 
c.Replace(1, m2)
```

### Working with Parent and Child Conversations
You can create child conversations and set parent conversations using the NewChild and WithParent methods. There are several use cases for this functionality:
- Internal monologue – create a child conversation from an existing one and modify the child conversation to better represent internal monologue, e.g. adding a prompt to the most recent message such as "is there enough information present in the conversation to answer this question?". Once the chatbot's next response is determined, it can be appended to the unmodified parent conversation.
- Summarizing a conversation without losing the original conversation's history – create a child conversation from an existing one, and append a message requesting a summary of the conversation. Send this to the chat completion endpoint. Once a response is received, create a new empty conversation using the original conversation as the parent. Append the summary message to the new conversation, and proceed with the conversation as normal. Whenever you need to refer to the original conversation, simply access it from the child conversation using its Parent() method. 

```go
// Create a child conversation. This automatically sets the child conversation's parent to the current conversation.
child := c.NewChild()

// Get the parent conversation
parent := child.Parent()

// Alternatively, create a new conversation and then set its parent conversation
c := conversation.New()
child := conversation.New(conversation.WithParent(c))
```

## License
This package is released under the MIT License. See [LICENSE](/LICENSE) for more information.