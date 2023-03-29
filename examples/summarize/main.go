package main

import (
	"context"
	"fmt"
	"github.com/bradfair/chat/conversation"
	"github.com/bradfair/chat/message"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
)

func main() {
	key := os.Getenv("OPENAI_APIKEY")
	if key == "" {
		log.Fatalln("Please set the OPENAI_APIKEY environment variable with your OpenAI API Key.")
	}

	originalConversation := conversation.New()
	originalConversation.WithMessages(
		message.New().WithRole("user").WithContent("Hi there. I'm wondering if you can help me with a computer problem."),
		message.New().WithRole("assistant").WithContent("Sure! Can you tell me a little more about what's going on?"),
		message.New().WithRole("user").WithContent("My computer is running really slow."),
		message.New().WithRole("assistant").WithContent("Okay, let's see what we can do. Have you tried restarting your computer?"),
		message.New().WithRole("user").WithContent("No, I haven't. How do I do that?"),
		message.New().WithRole("assistant").WithContent("I can walk you through it. Do you know what operating system you're using?"),
		message.New().WithRole("user").WithContent("It's a Windows computer."),
		message.New().WithRole("assistant").WithContent("Okay, let's try restarting your computer. First, click the Start button, then click the Power button, and then click Restart."),
		message.New().WithRole("user").WithContent("Okay. I'm restarting my computer now."),
		message.New().WithRole("assistant").WithContent("Great! Let me know when you're back."),
		message.New().WithRole("user").WithContent("I'm back."),
		message.New().WithRole("assistant").WithContent("Okay, let me know if you're still having problems."),
		message.New().WithRole("user").WithContent("It seems to be working now. Thanks for your help!"),
		message.New().WithRole("assistant").WithContent("Great! I'm glad I could help. Have a great day!"),
	)

	summaryPrompt := SummarizePrompt(originalConversation)

	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: ToChatCompletionMessages(summaryPrompt),
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}

func SummarizePrompt(c *conversation.Conversation) *conversation.Conversation {
	summary := c.NewChild().WithMessages(c.Messages()...)
	summary.Append(message.New().WithRole("system").WithContent("Without responding to any previous message, please briefly summarize the conversation so far."))
	return summary
}

// ToChatCompletionMessages converts a conversation to a slice of openai.ChatCompletionMessage.
func ToChatCompletionMessages(c *conversation.Conversation) []openai.ChatCompletionMessage {
	var messages []openai.ChatCompletionMessage
	for _, m := range c.Messages() {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    m.Role(),
			Content: m.Content(),
		})
	}
	return messages
}
