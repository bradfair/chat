# Chat Module
This module is designed to make it easy to represent and work with conversations in Golang. It consists of three packages: conversation, message, and tokenizer.

A conversation is an ordered collection of messages. A message is a string of text (the content) that is associated with
its sender (the role). Since OpenAI prices their API usage based on the combined number of tokens in a request/completion,
messages can make use of a tokenizer to convert their content into a list of integers (tokens). The Conversation struct
contains a helper method to count the total number of tokens in a conversation.

## Installation
To install this module, run the following command:

```sh
go get github.com/bradfair/chat
```

## Usage

Below is a brief overview of each package. For more detailed information, see the READMEs in each package's directory.

### Conversation Package
The [conversation package](conversation) provides a Conversation struct with relevant methods/functions to create and manipulate conversations.

### Message Package
The [message package](message) provides a Message struct with relevant methods/functions to create and manipulate messages.

### Tokenizer Package
The [tokenizer package](tokenizer) defines a Tokenizer interface and provides a TokenizerFunc type in order to satisfy the Tokenizer interface using custom functions. This allows you to use your own tokenizer function, or one of the small number of available open-source tokenizers.

## License
This module is licensed under the MIT License. See [LICENSE](LICENSE) for more information.