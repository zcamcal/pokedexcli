package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Welcome to the Pokedex!",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		t := scanner.Text()

		cleaned := cleanInput(t)
		fmt.Print("\n")
		v, ok := commands[cleaned[0]]
		if !ok {
			fmt.Println("command was not found")
			continue
		}

		v.callback()
	}
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	texts := strings.Fields(lowered)

	return texts
}
