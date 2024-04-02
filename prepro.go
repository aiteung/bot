package brain

import (
	"encoding/json"
	"os"
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
