package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/zcamcal/pokedexcli/internal/pokecache"
	"github.com/zcamcal/pokedexcli/internal/pokerepository"
)

type page struct {
	previous int
	next     int
}

type config struct {
	cache    *pokecache.Cache
	page     *page
	poketory pokerepository.PokeApi
}

type cliCommand struct {
	name        string
	description string
	callback    func(args []string, config *config) error
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

		"catch": {
			name:        "catch",
			description: "allows us to capture a pokemoon",
			callback:    commandCatch,
		},

		"explore": {
			name:        "explore",
			description: "explore more about the selected zone",
			callback:    commandExplore,
		},

		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	interval := 5 * time.Second
	cache := pokecache.NewCache(interval)
	poketory := pokerepository.NewPokeApi(interval)

	previous := -2
	next := 0
	actualPage := page{previous, next}
	config := &config{cache: &cache, poketory: poketory, page: &actualPage}

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

		err := v.callback(cleaned[1:], config)
		if err != nil {
			fmt.Printf("error from callback => %v", err)
		}
	}
}

func commandExplore(args []string, config *config) error {
	names, err := config.poketory.Encounters(args[0])
	if err != nil {
		return err
	}

	fmt.Println("Exploring pastoria-city-area...")
	fmt.Println("Found Pokemon:")
	for _, name := range names {
		fmt.Printf("- %v\n", name)
	}

	return nil
}

func commandCatch(args []string, config *config) error {
	pokemonName := args[0]
	experiencie, err := config.poketory.Pokemon(pokemonName)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonName)

	mostExperiencedPokemon := 1000
	chanceToCatch := rand.IntN(mostExperiencedPokemon)

	rate := mostExperiencedPokemon - experiencie

	if chanceToCatch <= rate {
		fmt.Printf("%v was caught!\n", pokemonName)
	} else {
		fmt.Printf("%v escaped!\n", pokemonName)
	}

	return nil
}

func commandMapBack(args []string, config *config) error {

	if config.page.next == 0 || config.page.previous < 0 {
		fmt.Println("you're on the first page")
		return nil
	}

	pokeIterations, err := config.poketory.Locations(20, config.page.previous)
	if err != nil {
		return err
	}
	config.page.next = config.page.previous + 1
	config.page.previous--

	for _, value := range pokeIterations {
		fmt.Println(value)
	}

	return nil
}

func commandMap(args []string, config *config) error {
	pokeIterations, err := config.poketory.Locations(20, config.page.next)
	if err != nil {
		return err
	}

	config.page.previous = config.page.next - 1
	config.page.next++

	for _, value := range pokeIterations {
		fmt.Println(value)
	}

	return nil
}

func commandHelp(args []string, config *config) error {
	fmt.Println("Welcome to the Pokedex!")
	return nil
}

func commandExit(args []string, config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	texts := strings.Fields(lowered)

	return texts
}

/*
func getLocationsPokemon[T any](path string, config config) (pokeIterator[T], error) {
	val, ok := config.cache.Get(path)
	if ok {
		var pokeResponse pokeIterator[T]

		if err := json.Unmarshal(val, &pokeResponse); err != nil {
			return pokeResponse, err
		}
	}

	res, err := http.Get(path)
	if err != nil {
		return pokeIterator[T]{}, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var pokeResponse pokeIterator[T]
	marshalled, err := json.Marshal(pokeResponse)
	if err != nil {
		return pokeIterator[T]{}, err
	}
	config.cache.Add(path, marshalled)

	if err := decoder.Decode(&pokeResponse); err != nil {
		return pokeIterator[T]{}, err
	}

	return pokeResponse, nil
}
*/
