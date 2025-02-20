package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Backend struct {
	ID   string `json:"id"`
	Port int    `json:"port"`
}

type Config struct {
	Backends []Backend `json:"backends"`
}

// Load reads the configuration from a JSON file and returns a Config instance
func Load(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config JSON: %w", err)
	}

	return &config, nil
}
