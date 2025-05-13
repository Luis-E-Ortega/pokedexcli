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
		if len(words) == 0 {
			continue
		}
		firstWord := words[0] // Capture the first word
		arg := ""
		// If there are multiple words then user is probably passing arguments
		if len(words) > 1 {
			arg = words[1]
		}
		// Check if command exists in registry
		if command, exists := commands[firstWord]; exists {
			err := command.callback(cfg, arg) // Callback passing the argument
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
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
