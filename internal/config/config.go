package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken  string
	DatabaseURL   string
	Port          string
	AllowedOrigin string
	BotDisabled   bool
}

func LoadConfig() *Config {
	// Load .env file if it exists, otherwise rely on system env vars
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")
	botDisabled := os.Getenv("BOT_DISABLED") == "true"
	if token == "" && !botDisabled {
		log.Println("DISCORD_TOKEN is not set; starting API only. Set BOT_DISABLED=true to silence this warning.")
		botDisabled = true
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL is not set; database-backed features will be disabled.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

	return &Config{
		DiscordToken:  token,
		DatabaseURL:   dbURL,
		Port:          port,
		AllowedOrigin: allowedOrigin,
		BotDisabled:   botDisabled,
	}
}
