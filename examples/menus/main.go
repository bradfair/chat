package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bradfair/chat/conversation"
	"github.com/bradfair/chat/message"
	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
	"regexp"
	"strconv"
)

type PromptDoer interface {
	Prompt() string
	WithInput(input string) PromptDoer
	WithParent(parent PromptDoer) PromptDoer
	Do() (output string, next PromptDoer)
}

type MenuItem struct {
	Title string
	PromptDoer
}

type Menu struct {
	Title  string
	input  string
	Parent PromptDoer
	Items  []MenuItem
}

func (m Menu) Prompt() string {
	var menuItems string
	for i, item := range m.Items {
		menuItems += fmt.Sprintf("%d. %s\n", i+1, item.Title)
	}
	return fmt.Sprintf("%s\n%s\nChoose a number: ", m.Title, menuItems)
}

func (m Menu) Do() (string, PromptDoer) {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(m.input)
	if match == "" {
		return fmt.Sprintf("invalid number: %.*s. Enter ONLY a number, and no other characters", 10, m.input), m
	}
	choice, err := strconv.Atoi(match)
	if err != nil || choice < 1 || choice > len(m.Items) {
		return fmt.Sprintf("invalid response: %q. try again.", match), m
	}
	return "", m.Items[choice-1]
}

func (m Menu) WithInput(input string) PromptDoer {
	m.input = input
	return m
}

func (m Menu) WithParent(parent PromptDoer) PromptDoer {
	m.Parent = parent
	return m
}

type Action struct {
	PromptText string
	input      string
	parent     PromptDoer
	DoFunc     func(input string, parent PromptDoer) (output string, next PromptDoer)
}

func (a Action) Prompt() string {
	return a.PromptText
}

func (a Action) Do() (string, PromptDoer) {
	return a.DoFunc(a.input, a.parent)
}

func (a Action) WithInput(input string) PromptDoer {
	a.input = input
	return a
}

func (a Action) WithParent(parent PromptDoer) PromptDoer {
	a.parent = parent
	return a
}

var output, goals, tasks, notes, rules string

func main() {
	openAiKey := os.Getenv("OPENAI_APIKEY")
	if openAiKey == "" {
		log.Fatalln("Please set the OPENAI_APIKEY environment variable with your OpenAI API Key.")
	}

	flag.StringVar(&goals, "goals", "", "The goals to accomplish.")
	flag.Parse()

	if goals == "" {
		log.Fatalln("Please set the goals to accomplish with the -goals flag.")
	}

	tasks = "1. Create a to-do list of all the tasks that need to be completed to accomplish the goal."

	originalConversation := conversation.New()
	originalConversation.Append(message.New().WithRole("system").WithContent(fmt.Sprintf("You are an AI talking to a menu in order to accomplish these goals:\n%s\n\nYou're currently working on the first task.", goals)))

	exitAction := &Action{DoFunc: func(input string, parent PromptDoer) (string, PromptDoer) { os.Exit(0); return "", nil }}

	mainMenu := &Menu{Title: "Main Menu"}

	tasksMenu := &Menu{Title: "Tasks Menu", Parent: mainMenu}
	notesMenu := &Menu{Title: "Notes Menu", Parent: mainMenu}

	viewTasksAction := &Action{DoFunc: func(input string, parent PromptDoer) (string, PromptDoer) { return tasks, parent }}
	editTasksAction := &Action{PromptText: "You can reprioritize, add, edit, or remove tasks here by replacing them. Replace tasks with: ", DoFunc: func(input string, parent PromptDoer) (string, PromptDoer) { tasks = input; return "", parent }}
	viewNotesAction := &Action{DoFunc: func(input string, parent PromptDoer) (string, PromptDoer) { return notes, parent }}
	editNotesAction := &Action{PromptText: "You can add, edit, or remove notes here by replacing them. Replace notes with: ", DoFunc: func(input string, parent PromptDoer) (string, PromptDoer) { notes = input; return "", parent }}

	tasksMenu.Items = []MenuItem{
		{"View Tasks", viewTasksAction},
		{"Edit Tasks", editTasksAction},
		{"Go Back", mainMenu},
	}

	notesMenu.Items = []MenuItem{
		{"View Notes", viewNotesAction},
		{"Edit Notes", editNotesAction},
		{"Go Back", mainMenu},
	}

	mainMenu.Items = []MenuItem{
		{"View/Edit Task List", tasksMenu},
		{"View/Edit Notes", notesMenu},
		{"Exit", exitAction},
	}

	var currentAction PromptDoer = mainMenu
	for {
		var input string
		prompt := currentAction.Prompt()
		if prompt != "" {
			originalConversation.Append(message.New().WithRole("user").WithContent(output + "\n" + prompt))
			fmt.Print(prompt)
			input = ThinkAndRespond(openAiKey, originalConversation)
			originalConversation.Append(message.New().WithRole("assistant").WithContent(input))
			color.Set(color.FgBlue)
			fmt.Println(input)
			color.Unset()
		}
		previousAction := currentAction
		output, currentAction = currentAction.WithInput(input).Do()
		if &currentAction != &previousAction {
			currentAction = currentAction.WithParent(previousAction)
		}
		if output != "" {
			color.Set(color.FgGreen)
			fmt.Println(output)
			color.Unset()
		}
	}
}

func ThinkAndRespond(openAiKey string, originalConversation *conversation.Conversation) string {
	trimmedConversation := originalConversation.NewChild()
	tailLength := 20
	convLen := len(originalConversation.Messages()[1:])

	for i := convLen - 1; i >= 0; i-- {
		if i < convLen-tailLength {
			break
		}
		trimmedConversation.Prepend(originalConversation.Messages()[i+1])
	}
	trimmedConversation.Prepend(originalConversation.Messages()[0])
	// Uncomment to see the full conversation that is being sent to OpenAI for each request. It's basically the initial prompt + the last 20 messages.
	//color.Set(color.FgYellow)
	//fmt.Println(trimmedConversation.Messages().Transcript())
	//color.Unset()
	assistantResponse, err := getCompletion(openAiKey, trimmedConversation)
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
			// GPT3Dot5Turbo is faster and cheaper than GPT4. It's sufficient for this demo, and stays on track well enough.
			//Model: openai.GPT4,
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
