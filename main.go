package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration fields.
type Config struct {
	ServerPort  string `yaml:"server_port"`
	DatabaseURL string `yaml:"database_url"`
	LogLevel    string `yaml:"log_level"`
}

// Validate checks that all required configuration fields are present and non-empty.
func (c *Config) Validate() error {
	checks := []struct {
		name  string
		value string
	}{
		{"server_port", c.ServerPort},
		{"database_url", c.DatabaseURL},
		{"log_level", c.LogLevel},
	}
	for _, ch := range checks {
		if ch.value == "" {
			return fmt.Errorf("missing required config field: %s", ch.name)
		}
	}
	return nil
}

// LoadConfig reads and parses a YAML configuration file, returning a validated Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %s: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return &cfg, nil
}

func main() {
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Config loaded successfully:\n")
	fmt.Printf("  server_port:  %s\n", cfg.ServerPort)
	fmt.Printf("  database_url: %s\n", cfg.DatabaseURL)
	fmt.Printf("  log_level:    %s\n", cfg.LogLevel)
}
