package config

import (
	"fmt"
	"os"
)

// Config sadrži postavke potrebne za pokretanje servisa.
type Config struct {
	DatabaseURL string
	Port        string
}

// Load čita postavke iz varijabli okoline i vraća grešku ako neka nedostaje.
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
