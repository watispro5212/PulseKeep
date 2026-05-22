package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken string
	DatabaseURL  string
	Port         string
}

func LoadConfig() *Config {
	// Load .env file if it exists, otherwise rely on system env vars
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN is required")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DiscordToken: token,
		DatabaseURL:  dbURL,
		Port:         port,
	}
}
