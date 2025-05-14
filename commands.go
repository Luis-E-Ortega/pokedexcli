package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/Luis-E-Ortega/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

type config struct {
	Next          *string
	Prev          *string
	cache         *pokecache.Cache
	caughtPokemon map[string]Pokemon
}

type locationAreaResponse struct {
	Next    *string            `json:"next"`
	Prev    *string            `json:"previous"`
	Results []locationAreaInfo `json:"results"`
}

type locationAreaInfo struct {
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name           string        `json:"name"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
}

type PokemonStat struct {
	BaseStat int              `json:"base_stat"`
	Effort   int              `json:"effort"`
	Stat     NamedAPIResource `json:"stat"`
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonType struct {
	Slot int              `json:"slot"`
	Type NamedAPIResource `json:"type"`
}

// Function that returns the map of commands
func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "See the list of all pokemon at a given location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Adds pokemon to Pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect details about caught pokemon in pokedex",
			callback:    commandInspect,
		},
	}
}

func commandExit(cfg *config, name string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, name string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	// Just a note for memory efficiency, storing the map in memory on each command call isn't as good
	// as making the commands a package level variable. It just works fine in this small program
	for _, c := range getCommands() {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMap(cfg *config, name string) error {
	url := "https://pokeapi.co/api/v2/location-area"
	var body []byte
	var err error

	// If a "next" URL already exists, use that instead
	if cfg.Next != nil {
		url = *cfg.Next
	}
	// Check if response is in cache
	if cacheData, found := cfg.cache.Get(url); found {
		fmt.Println("Cache found!")
		body = cacheData
	} else {
		// Making the get request to pull API location data
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		// Read the body of the data
		body, err = io.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			return err
		}

		if res.StatusCode > 299 {
			return fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
		}

		// Add to cache
		cfg.cache.Add(url, body)
	}

	locations := locationAreaResponse{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return err
	}

	cfg.Next = locations.Next
	cfg.Prev = locations.Prev

	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapb(cfg *config, name string) error {
	if cfg.Prev != nil {
		url := *cfg.Prev
		var body []byte
		var err error

		// Check if response is in cache
		if cachedData, found := cfg.cache.Get(url); found {
			fmt.Println("Using cached data!")
			body = cachedData
		} else {
			res, err := http.Get(url)
			if err != nil {
				return err
			}

			body, err = io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				return err
			}

			if res.StatusCode > 299 {
				return fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
			}

			// Add to cache
			cfg.cache.Add(url, body)
		}

		// Process the data (cached or freshly fetched)
		locations := locationAreaResponse{}
		err = json.Unmarshal(body, &locations)
		if err != nil {
			return err
		}

		for _, area := range locations.Results {
			fmt.Println(area.Name)
		}

		cfg.Next = locations.Next
		cfg.Prev = locations.Prev
	} else {
		fmt.Println("you're on the first page")
	}

	return nil
}

func commandExplore(cfg *config, name string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + name + "/"

	var body []byte
	var err error

	// Check if response is in cache
	if cacheData, found := cfg.cache.Get(url); found {
		fmt.Println("Cache found!")
		body = cacheData
	} else {
		// Making the get request to pull API location data
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		// Read the body of the data
		body, err = io.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			return err
		}

		if res.StatusCode > 299 {
			return fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
		}

		// Add to cache
		cfg.cache.Add(url, body)
	}

	location := locationAreaInfo{}
	err = json.Unmarshal(body, &location)
	if err != nil {
		return err
	}

	if name != "" {
		// To print the list of pokemon from location
		fmt.Printf("Exploring %s...\n", name)
		fmt.Printf("Found Pokemon:\n")
	}

	for _, area := range location.PokemonEncounters {
		fmt.Printf("- %s\n", area.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, name string) error {
	if name == "" {
		return fmt.Errorf("requires a valid pokemon name")
	}
	// To format the url and name correctly
	lowerName := strings.ToLower(name)
	url := "https://pokeapi.co/api/v2/pokemon/" + lowerName + "/"

	var body []byte
	var err error

	// If the data is already found in cache
	if cachedData, found := cfg.cache.Get(url); found {
		body = cachedData
	} else {
		// Making the get request to pull API pokemon data
		res, err := http.Get(url)
		if err != nil {
			return err
		}

		// Ready the body of the data
		body, err = io.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			return err
		}
		if res.StatusCode == 404 {
			return fmt.Errorf("Pokemon '%s' not found", lowerName)
		} else if res.StatusCode > 299 {
			return fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
		}

		// Add to cache
		cfg.cache.Add(url, body)
	}

	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return err
	}

	randInt := rand.Intn(100)
	caught := false

	fmt.Printf("Throwing a Pokeball at %s...\n", lowerName)
	if pokemon.BaseExperience < 40 {
		if randInt <= 75 {
			caught = true
		}
	} else if pokemon.BaseExperience <= 150 {
		if randInt <= 35 {
			caught = true
		}
	} else if pokemon.BaseExperience <= 300 {
		if randInt <= 20 {
			caught = true
		}
	} else {
		if randInt <= 10 {
			caught = true
		}
	}

	if caught {
		fmt.Printf("%s was caught!\n", lowerName)
		cfg.caughtPokemon[lowerName] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", lowerName)
	}

	return nil
}

func commandInspect(cfg *config, name string) error {
	lowerName := strings.ToLower(name)

	if name == "" {
		return fmt.Errorf("you must provide a pokemon name")
	}

	if pokeData, ok := cfg.caughtPokemon[lowerName]; ok {
		// Print simple output together with newlines
		fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nStats:\n", pokeData.Name, pokeData.Height, pokeData.Weight)
		// Loop through stats
		for _, stat := range pokeData.Stats {
			fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}

		fmt.Println("Types:")
		// Loop through types
		for _, t := range pokeData.Types {
			fmt.Printf("  - %s\n", t.Type.Name)
		}
	} else {
		fmt.Println("you have not caught that pokemon")
	}
	return nil
}
