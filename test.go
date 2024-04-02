package brain

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	// Example of using loadConfig and normalizeSentence
	configPath := "config.json"
	config, err := loadConfig(configPath)
	println(config)
	if err != nil {
		panic(err)
	}

	// Example usage of normalizeSentence
	sentence := "Example sentence to normalize!"
	normalized := normalizeSentence(sentence)
	println(normalized)

	// Example usage
	tokenizer, err := LoadTokenizer("/path/to", "tokenizer.json")
	if err != nil {
		// Handle error
		panic(err)
	}

	// Use the tokenizer
	// For example, print the tokenizer data
	fmt.Printf("Loaded tokenizer: %+v\n", tokenizer)

	fileName := "data"
	stemmer, punctReEscape, unknowns, path := setConfig(fileName)
	fmt.Println(stemmer, punctReEscape, unknowns, path)

	stemmer = NewStemmer()
	words := []string{"berjalan", "menggunakan", "pengetahuan", "terapkan"}

	for _, word := range words {
		stemmedWord := stemmer.Stem(word)
		fmt.Printf("Original: %-12s Stemmed: %-12s\n", word, stemmedWord)
	}
}
