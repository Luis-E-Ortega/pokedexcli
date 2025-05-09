package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Luis-E-Ortega/pokedexcli/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin) // Initialize scanner
	commands := getCommands()             // Get the commands map once at the start
	// Create an instance of the config struct and initialize cache
	cfg := &config{
		cache: pokecache.NewCache(5 * time.Minute),
	}

	// Infinite loop while program is running to register user input
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()          // Scans for last input
		input := scanner.Text() // Save input text into variable

		words := cleanInput(input) // Split input into words

		// Check if there is at least one word
		if len(words) > 0 {
			firstWord := words[0] // Capture first word
			// Check if command exists in registry
			if command, exists := commands[firstWord]; exists {
				// Call its callback
				err := command.callback(cfg)
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				// Command not found
				fmt.Println("Unknown command")
			}
		} else {
			// Empty Input
			continue
		}

	}
}

func cleanInput(text string) []string {
	// String cleanup for space and lowercase
	lowercased := strings.ToLower(text)

	// Splitting lowercased strings into slice of words
	cleanedStrings := strings.Fields(lowercased)

	return cleanedStrings
}
