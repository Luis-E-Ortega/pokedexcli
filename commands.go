package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Luis-E-Ortega/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next  *string
	Prev  *string
	cache *pokecache.Cache
}

type locationAreaResponse struct {
	Next    *string            `json:"next"`
	Prev    *string            `json:"previous"`
	Results []locationAreaInfo `json:"results"`
}

type locationAreaInfo struct {
	Name string `json:"name"`
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
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	// Just a note for memory efficiency, storing the map in memory on each command call isn't as good
	// as making the commands a package level variable. It just works fine in this small program
	for _, c := range getCommands() {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMap(cfg *config) error {
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

func commandMapb(cfg *config) error {
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
