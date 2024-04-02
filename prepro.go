package brain

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// loadConfig reads the configuration from a JSON file.
func loadConfig(path string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// normalizeSentence cleans and normalizes sentences.
func normalizeSentence(sentence string, punctRe *regexp.Regexp) string {
	sentence = punctRe.ReplaceAllString(strings.ToLower(sentence), "")
	sentence = strings.ReplaceAll(sentence, "iiteung", "")
	sentence = strings.ReplaceAll(sentence, "iteung", "")
	sentence = strings.ReplaceAll(sentence, "teung", "")
	sentence = strings.ReplaceAll(sentence, "\n", "")
	// Add more replacements and regex substitutions as needed
	sentence = strings.TrimSpace(sentence)
	return sentence
}

type Tokenizer struct {
	// Example fields
	Tokens map[string]int
	// Add other fields as necessary
}

// LoadTokenizer loads a tokenizer from a file in JSON format.
func LoadTokenizer(basePath, tokenizerPath string) (*Tokenizer, error) {
	filePath := filepath.Join(basePath, tokenizerPath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tokenizer Tokenizer
	if err := json.Unmarshal(data, &tokenizer); err != nil {
		return nil, err
	}

	return &tokenizer, nil
}

type Stemmer struct{}

// stringJoin joins strings with a separator, similar to strings.Join
func stringJoin(separator string, elements []string) string {
	result := ""
	for i, el := range elements {
		if i > 0 {
			result += separator
		}
		result += el
	}
	return result
}

// Stem reduces a word to its base form
func (s *Stemmer) Stem(word string) string {
	// Define common prefixes and suffixes
	prefixes := []string{"ber", "ter", "meng", "peng"}
	suffixes := []string{"kan", "i", "an"}

	// Compile regular expressions for prefixes and suffixes
	prefixRe := regexp.MustCompile("^(?i)(" + stringJoin("|", prefixes) + ")")
	suffixRe := regexp.MustCompile("(?i)(" + stringJoin("|", suffixes) + ")$")

	// Remove common prefixes
	word = prefixRe.ReplaceAllString(word, "")

	// Remove common suffixes
	word = suffixRe.ReplaceAllString(word, "")
	return word
}

func NewStemmer() *Stemmer {
	return &Stemmer{}
}

func setConfig(fileName string) (*Stemmer, *regexp.Regexp, []string, string) {
	stemmer := NewStemmer()
	punctReEscape := regexp.MustCompile(`[!"#$%&'()*+,\-./:;<=>?@[\\\]^_` + "`{|}~]")
	unknowns := []string{"gak paham", "kurang ngerti", "I don't know"}
	path := filepath.Join(fileName, "/")

	return stemmer, punctReEscape, unknowns, path
}
