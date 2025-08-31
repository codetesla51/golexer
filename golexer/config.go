// golexer/config.go
package golexer

import (
	"encoding/json"
	"os"
)

type Config struct {
	AdditionalKeywords    map[string]string `json:"additionalKeywords"`
	AdditionalOperators   map[string]string `json:"additionalOperators"`
	AdditionalPunctuation map[string]string `json:"additionalPunctuation"`
}

func (c *Config) MergeWithDefaults() {
	for keyword, tokenType := range c.AdditionalKeywords {
		keywords[keyword] = TokenType(tokenType)
	}

	for op, tokenType := range c.AdditionalOperators {
		operators = append(operators, Operator{
			Single:     op,
			SingleType: TokenType(tokenType),
		})
	}

	for char, tokenType := range c.AdditionalPunctuation {
		if len(char) == 1 {
			singleCharTokens[rune(char[0])] = TokenType(tokenType)
		}
	}
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	return &config, err
}
