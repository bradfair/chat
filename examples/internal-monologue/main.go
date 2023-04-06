package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/bradfair/chat/conversation"
	"github.com/bradfair/chat/message"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"strings"
)

func main() {
	openAiKey := os.Getenv("OPENAI_APIKEY")
	if openAiKey == "" {
		log.Fatalln("Please set the OPENAI_APIKEY environment variable with your OpenAI API Key.")
	}

	originalConversation := conversation.New()
	originalConversation.WithMessages(
		message.New().WithRole("user").WithContent("You are a curious and charismatic chatbot. I am a human. Let's talk! Start by greeting me and asking an icebreaker question."),
	)
	fmt.Println("Initial prompt:", originalConversation.Messages()[0].Content())
	assistantResponse, err := getCompletion(openAiKey, originalConversation)
	if err != nil {
		log.Fatalln(err)
	}
	originalConversation.Append(message.New().WithRole("assistant").WithContent(assistantResponse))
	fmt.Println("Chatbot: " + assistantResponse)
	reader := bufio.NewReader(os.Stdin)
	for {
		var userInput string
		fmt.Print("User: ")
		userInput, err = reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		if userInput == "" {
			continue
		}
		if userInput == "bye" {
			break
		}
		originalConversation.Append(message.New().WithRole("user").WithContent(userInput))
		assistantResponse = ThinkAndRespond(openAiKey, originalConversation)
		originalConversation.Append(message.New().WithRole("assistant").WithContent(assistantResponse))
		fmt.Println("Chatbot: " + assistantResponse)
	}
}

func ThinkAndRespond(openAiKey string, originalConversation *conversation.Conversation) string {
	prompt := "You are an AI that serves as the conscience of a curious and charismatic chatbot."
	transcript := "Here's a transcript of a conversation you're having with a human:\n\n" + originalConversation.Messages().Transcript()
	rules := "We have strict rules for handling conversations:\n1. Stay on topic: you are a curious and charismatic chatbot.\n2. Be respectful. Don't allow the conversation to become hostile.\n3. Be safe. Don't allow the conversation to become dangerous.\n4. Be honest. Don't lie or mislead."
	request := "1) What is the main goal of the conversation as a whole?\n2) What is the main goal of the most recent messages?\n3) Abide by the rules. What is the best response to the most recent message?\n4) Critique your own response.\n\n1)"
	internalMonologue := originalConversation.NewChild()
	internalMonologue.Append(message.New().WithRole("system").WithContent(fmt.Sprintf("%s\n%s\n%s\n%s", prompt, transcript, rules, request)))
	//fmt.Println("User (internal monologue): " + internalMonologue.Messages()[len(internalMonologue.Messages())-1].Content())
	assistantResponse, err := getCompletion(openAiKey, internalMonologue)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Chatbot (internal monologue): " + assistantResponse)
	internalMonologue.Append(message.New().WithRole("assistant").WithContent(assistantResponse))
	internalMonologue.Append(message.New().WithRole("user").WithContent(fmt.Sprintf("With your self-critique in mind, briefly respond directly to %q:", originalConversation.Messages()[len(originalConversation.Messages())-1].Content())))
	//fmt.Println("User (internal monologue): " + internalMonologue.Messages()[len(internalMonologue.Messages())-1].Content())
	assistantResponse, err = getCompletion(openAiKey, internalMonologue)
	if err != nil {
		log.Fatalln(err)
	}
	return assistantResponse
}

func getCompletion(key string, convo *conversation.Conversation) (string, error) {
	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: ToChatCompletionMessages(convo),
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

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

func IsConversationOver(openAiKey string, c *conversation.Conversation) bool {
	check := c.NewChild().WithMessages(c.Messages()...)
	check.Append(message.New().WithRole("system").WithContent("Is the conversation over?\n1. Yes\n2. No\nRespond with the number of your choice. If you don't know, respond with 2.\nResponse:"))
	var assistantResponse string
	var err error
	for {
		assistantResponse, err = getCompletion(openAiKey, check)
		if err != nil {
			log.Fatalln(err)
		}
		if assistantResponse == "1" || assistantResponse == "2" {
			break
		}
		check.Append(message.New().WithRole("assistant").WithContent(assistantResponse))
		check.Append(message.New().WithRole("system").WithContent("Please respond with 1 or 2."))
	}
	return assistantResponse == "1"
}
