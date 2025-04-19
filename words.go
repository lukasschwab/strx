package main

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
)

// randomWords selects two random words from /usr/share/dict/words. It panics if
// anything goes wrong; brittle!
func randomWords() (string, string) {
	file, err := os.Open("/usr/share/dict/words")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read all words into a slice
	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// TODO: more thoroughly consider character sets.
		if len(scanner.Text()) > 5 {
			continue
		} else if strings.Contains(scanner.Text(), " ") {
			continue
		}

		words = append(words, strings.ToLower(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	if len(words) < 2 {
		panic("Not enough words in the dictionary")
	}
	word1 := words[rand.Intn(len(words))]
	word2 := words[rand.Intn(len(words))]

	return word1, word2
}
