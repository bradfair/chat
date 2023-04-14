package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	mainMenu := new(Action)
	submenu1 := new(Action)
	submenu2 := new(Action)
	subsubmenu1 := new(Action)
	action1 := new(Action)
	action2 := new(Action)

	mainMenu.Title = "Main Menu"
	mainMenu.Prompt = fmt.Sprintf("%s\n1. sub-menu 1\n2. sub-menu 2\n3. Exit\n\nChoice: ", mainMenu.Title)
	mainMenu.Do = func() (*Action, error) {
		switch mainMenu.Input {
		case "1":
			return submenu1, nil
		case "2":
			return submenu2, nil
		case "3":
			os.Exit(0)
		}
		return mainMenu, fmt.Errorf("invalid input: %s", mainMenu.Input)
	}

	submenu1.Title = "sub-menu 1"
	submenu1.Prompt = fmt.Sprintf("%s\n1. sub-sub-menu 1\n2. action that prompts for input\n3. action that does not prompt for input\n4. Go back\n\nChoice: ", submenu1.Title)
	submenu1.Parent = mainMenu
	submenu1.Do = func() (*Action, error) {
		switch submenu1.Input {
		case "1":
			return subsubmenu1, nil
		case "2":
			action1.Parent = submenu1
			return action1, nil
		case "3":
			action2.Parent = submenu1
			return action2, nil
		case "4":
			return mainMenu, nil
		}
		return submenu1, fmt.Errorf("invalid input: %s", submenu1.Input)
	}

	submenu2.Title = "sub-menu 2"
	submenu2.Prompt = fmt.Sprintf("%s\n1. action that prompts for input\n2. action that does not prompt for input\n3. Go back\n\nChoice: ", submenu2.Title)
	submenu2.Parent = mainMenu
	submenu2.Do = func() (*Action, error) {
		switch submenu2.Input {
		case "1":
			action1.Parent = submenu2
			return action1, nil
		case "2":
			action2.Parent = submenu2
			return action2, nil
		case "3":
			return mainMenu, nil
		}
		return submenu2, fmt.Errorf("invalid input: %s", submenu2.Input)
	}

	subsubmenu1.Title = "sub-sub-menu 1"
	subsubmenu1.Prompt = fmt.Sprintf("%s\n1. Go back\n\nChoice: ", subsubmenu1.Title)
	subsubmenu1.Parent = submenu1
	subsubmenu1.Do = func() (*Action, error) {
		switch subsubmenu1.Input {
		case "1":
			return submenu1, nil
		}
		return subsubmenu1, fmt.Errorf("invalid input: %s", subsubmenu1.Input)
	}

	action1.Title = "action that prompts for input"
	action1.Prompt = fmt.Sprintf("%s\nWhat would you like to say? ", action1.Title)
	action1.Do = func() (*Action, error) {
		fmt.Printf("You said: %s\n", action1.Input)
		return action1.Parent, nil
	}

	action2.Title = "action that does not prompt for input"
	action2.Do = func() (*Action, error) {
		fmt.Println("Hello, world!")
		return action2.Parent, nil
	}

	// Start at the main menu
	currentAction := mainMenu
	for {
		if currentAction.Prompt != "" {
			fmt.Print(currentAction.Prompt)
			var input string
			fmt.Scanln(&input)
			currentAction.Input = strings.TrimSpace(input)
		}
		nextAction, err := currentAction.Do()
		if err != nil {
			fmt.Println(err)
		}
		currentAction = nextAction
	}

}

type Action struct {
	Title  string
	Prompt string
	Input  string
	Do     func() (*Action, error)
	Parent *Action
}
