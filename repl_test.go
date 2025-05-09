package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "    hello world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  chEcKING  tHIs    NoW   ",
			expected: []string{"checking", "this", "now"},
		},
		{
			input:    "POKEmoN   are    fun   and interESTING",
			expected: []string{"pokemon", "are", "fun", "and", "interesting"},
		},
		{
			input:    "   joy",
			expected: []string{"joy"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("For input '%s': Expected length %d, got %d", c.input, len(c.expected), len(actual))
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("For input '%s': Expected word at index %d to be '%s', got '%s'", c.input, i, expectedWord, word)
			}
		}
	}
}
