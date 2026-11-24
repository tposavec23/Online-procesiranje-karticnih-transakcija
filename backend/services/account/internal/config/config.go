package config

import (
	"fmt"
	"os"
)

// Config contains the settings required to run the service.
type Config struct {
	DatabaseURL string
	Port        string
}

// Load reads the settings from environment variables and returns an error if any of them is missing.
func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return Config{}, fmt.Errorf("PORT is not set")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is not set")
	}

	return Config{
		DatabaseURL: databaseURL,
		Port:        port,
	}, nil
}
