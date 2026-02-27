package orchestrator

import (
	"encoding/json"
	"os"
)

// Config holds pathway and layer configuration.
type Config struct {
	Pathway string   `json:"pathway"`
	Layers  []string `json:"layers,omitempty"`
}

// LoadConfig reads config from a JSON file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// LoadConfigOrDefault returns config from path, or default Pathway A.
func LoadConfigOrDefault(path string) *Config {
	cfg, err := LoadConfig(path)
	if err != nil {
		return &Config{Pathway: "a"}
	}
	return cfg
}
