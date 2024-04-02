package brain

import (
	"regexp"
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
	punctRe := regexp.MustCompile(`[\!"#$%&'()*+,./:;<=>?@\[\\\]^_` + "`{|}~]")
	sentence := "Example sentence to normalize!"
	normalized := normalizeSentence(sentence, punctRe)
	println(normalized)
}
