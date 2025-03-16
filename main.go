package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type pokeIterator[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}

type pokeLocation struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type config struct {
	next     string
	previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *config) error
}

func main() {
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Welcome to the Pokedex!",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "help us gathering all maps from some configured location",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "help us gatherin the same page we gather before",
			callback:    commandMapBack,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	config := &config{}

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

		err := v.callback(config)
		if err != nil {
			fmt.Printf("error from callback => %v", err)
		}
	}
}

func getLocationsPokemon(path string) (pokeIterator[pokeLocation], error) {
	res, err := http.Get(path)
	if err != nil {
		return pokeIterator[pokeLocation]{}, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var pokeResponse pokeIterator[pokeLocation]
	if err := decoder.Decode(&pokeResponse); err != nil {
		return pokeIterator[pokeLocation]{}, err
	}

	return pokeResponse, nil
}

func commandMapBack(config *config) error {

	if config.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	pokeIterations, err := getLocationsPokemon(config.previous)
	if err != nil {
		return err
	}

	if pokeIterations.Previous != "" {
		config.previous = pokeIterations.Previous
	} else {
		config.previous = ""
	}

	if pokeIterations.Next != "" {
		config.next = pokeIterations.Next
	} else {
		config.next = "https://pokeapi.co/api/v2/location-area"
	}

	for _, value := range pokeIterations.Results {
		fmt.Println(value.Name)
	}

	return nil
}

func commandMap(config *config) error {
	path := "https://pokeapi.co/api/v2/location-area"

	if config.next != "" {
		path = config.next
	}

	pokeIterations, err := getLocationsPokemon(path)
	if err != nil {
		return err
	}

	if pokeIterations.Previous != "" {
		config.previous = pokeIterations.Previous
	}

	if pokeIterations.Next != "" {
		config.next = pokeIterations.Next
	}

	for _, value := range pokeIterations.Results {
		fmt.Println(value.Name)
	}

	return nil
}

func commandHelp(config *config) error {
	fmt.Println("Welcome to the Pokedex!")
	return nil
}

func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	texts := strings.Fields(lowered)

	return texts
}
