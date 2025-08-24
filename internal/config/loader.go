package config

import (
	"encoding/json"
	"os"
	"strings"
)

// Load parses a Config either from a raw JSON string
// or from a JSON file at the given path.
func Load(input string) (*Config, error) {
	var data []byte
	if strings.HasPrefix(strings.TrimSpace(input), "{") {
		// Raw JSON string
		data = []byte(input)
	} else {
		// Treat as file path
		b, err := os.ReadFile(input)
		if err != nil {
			return nil, err
		}
		data = b
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
