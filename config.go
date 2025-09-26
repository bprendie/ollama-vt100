// config/config.go

package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = "config.json"

// AppConfig holds the application's startup configuration.
type AppConfig struct {
	OllamaURL string `json:"ollama_url"`
}

// LoadConfig reads the configuration file. If the file doesn't exist,
// it creates a default one.
func LoadConfig() (AppConfig, error) {
	var cfg AppConfig

	file, err := os.Open(configFileName)
	if err != nil {
		// If file doesn't exist, create it with defaults
		if os.IsNotExist(err) {
			fmt.Printf("Config file not found. Creating a default '%s'.\n", configFileName)
			return createDefaultConfig()
		}
		// For any other error, fail
		return AppConfig{}, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return AppConfig{}, fmt.Errorf("could not decode config file: %w", err)
	}

	return cfg, nil
}

// createDefaultConfig writes a new config file with default values and returns it.
func createDefaultConfig() (AppConfig, error) {
	defaultConfig := AppConfig{
		OllamaURL: "http://localhost:11434",
	}

	file, err := os.Create(configFileName)
	if err != nil {
		return AppConfig{}, fmt.Errorf("could not create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Make it pretty
	if err := encoder.Encode(defaultConfig); err != nil {
		return AppConfig{}, fmt.Errorf("could not write to config file: %w", err)
	}

	return defaultConfig, nil
}
