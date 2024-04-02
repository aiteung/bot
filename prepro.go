package brain

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
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
func normalizeSentence(sentence string) string {
	punctRe := regexp.MustCompile(`[\!"#$%&'()*+,./:;<=>?@\[\\\]^_` + "`{|}~]")
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
	WordIndex map[string]int // WordIndex is a field of type map[string]int

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

// Let's assume these are pre-defined elsewhere in your Go application
var (
	stemmer         *Stemmer // Your stemmer implementation
	unknowns        = []string{"gak paham", "kurang ngerti", "I don't know"}
	tokenizer       *Tokenizer // Your tokenizer implementation
	maxlenAnswers   int
	encoderModel    *Model // Encapsulate your encoder model interaction in this struct
	decoderModel    *Model // Encapsulate your decoder model interaction in this struct
	maxlenQuestions int
)

type Model struct {
	// Add model fields and methods
}

// Predict method for Model
func (m *Model) Predict(inputs []int) ([]float64, error) {
	// Prediction logic here
	return []float64{}, nil // placeholder return
}

// Tokenize method for Tokenizer
func (t *Tokenizer) Tokenize(sentence string) []int {
	// Tokenization logic here
	return []int{} // placeholder return
}

// Chat function takes the input and converses based on the trained models
func Chat(inputValue string) (string, string) {
	// Preprocess and stem the input
	normalizedInput := stemmer.Stem(normalizeSentence(inputValue))
	tokens := tokenizer.Tokenize(normalizedInput)

	// Predict using encoder model
	statesValues, _ := encoderModel.Predict(tokens)

	emptyTargetSeq := make([]float64, 1) // Simulating numpy.zeros((1,1))
	emptyTargetSeq[0] = float64(tokenizer.WordIndex("start"))

	stopCondition := false
	var decodedTranslation strings.Builder
	status := "false"

	for !stopCondition {
		decOutputs, h, c := decoderModel.Predict(emptyTargetSeq, statesValues)

		sampledWordIndex := argmax(decOutputs)
		if decOutputs[sampledWordIndex] < 0.1 {
			randomIndex := rand.Intn(len(unknowns))
			decodedTranslation.WriteString(unknowns[randomIndex])
			break
		}

		sampledWord, exists := tokenizer.IndexToWord(sampledWordIndex)
		if !exists || sampledWord == "end" || decodedTranslation.Len() > maxlenAnswers {
			stopCondition = true
		} else {
			decodedTranslation.WriteString(" " + sampledWord)
		}

		emptyTargetSeq[0] = float64(sampledWordIndex)
		statesValues = []float64{h, c}
		status = "true"
	}

	return decodedTranslation.String(), status
}
