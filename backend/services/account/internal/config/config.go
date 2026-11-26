package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

// Config contains the settings required to run the service.
type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
}

// Load reads the settings from environment variables and returns an error if any required one is missing.
// JWT_SECRET is optional: if it is not set, a random secret is generated on every start.
func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return Config{}, fmt.Errorf("PORT is not set")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		generated, err := generateSecret()
		if err != nil {
			return Config{}, fmt.Errorf("generate JWT secret: %w", err)
		}
		jwtSecret = generated
		log.Println("JWT_SECRET is not set; using a random secret, tokens will be invalid after restart")
	}

	return Config{
		DatabaseURL: databaseURL,
		Port:        port,
		JWTSecret:   jwtSecret,
	}, nil
}

func generateSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
