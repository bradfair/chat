package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	return fmt.Sprintf("%s\n%s\nChoice: ", m.Title, menuItems)
}

func (m Menu) Do() (string, PromptDoer) {
	choice, err := strconv.Atoi(m.input)
	if err != nil || choice < 1 || choice > len(m.Items) {
		return fmt.Sprintf("invalid input: %s", m.input), m
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

func main() {
	mainMenu := &Menu{Title: "Main Menu"}
	submenu1 := &Menu{Title: "sub-menu 1"}
	submenu2 := &Menu{Title: "sub-menu 2"}

	subsubmenu1 := &Menu{Title: "sub-sub-menu 1"}

	action1 := &Action{
		PromptText: "What would you like to say? ",
		DoFunc: func(input string, parent PromptDoer) (string, PromptDoer) {
			return fmt.Sprintf("you said: %s", input), parent
		},
	}

	action2 := &Action{DoFunc: func(input string, parent PromptDoer) (string, PromptDoer) { return "Hello, World!", parent }}

	exitAction := &Action{DoFunc: func(input string, parent PromptDoer) (string, PromptDoer) { os.Exit(0); return "", nil }}

	mainMenu.Items = []MenuItem{
		{"sub-menu 1", submenu1},
		{"sub-menu 2", submenu2},
		{"Exit", exitAction},
	}

	submenu1.Items = []MenuItem{
		{"sub-sub-menu 1", subsubmenu1},
		{"action that prompts for input", action1.WithParent(submenu1)},
		{"action that does not prompt for input", action2.WithParent(submenu1)},
		{"Go back", mainMenu},
	}
	submenu1.Parent = mainMenu

	submenu2.Items = []MenuItem{
		{"action that prompts for input", action1.WithParent(submenu2)},
		{"action that does not prompt for input", action2.WithParent(submenu2)},
		{"Go back", mainMenu},
	}
	submenu2.Parent = mainMenu

	subsubmenu1.Items = []MenuItem{
		{"Go back", submenu1},
	}
	subsubmenu1.Parent = submenu1

	var currentAction PromptDoer = mainMenu

	scanner := bufio.NewScanner(os.Stdin)
	for {
		var input, output string
		prompt := currentAction.Prompt()
		if prompt != "" {
			fmt.Print(prompt)
			scanner.Scan()
			input = scanner.Text()
			input = strings.TrimSpace(input)
			if input == "" {
				continue
			}
		}
		output, currentAction = currentAction.WithInput(input).Do()
		fmt.Println(output)
	}
}
