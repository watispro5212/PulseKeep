package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/watispro/pulsekeep/internal/api"
	"github.com/watispro/pulsekeep/internal/bot"
	"github.com/watispro/pulsekeep/internal/cache"
	"github.com/watispro/pulsekeep/internal/config"
	"github.com/watispro/pulsekeep/internal/db"
)

func main() {
	cfg := config.LoadConfig()
	database := db.Connect(cfg.DatabaseURL)
	defer database.Close()
	memCache := cache.New()
	_ = memCache

	discordBot := bot.New(cfg.DiscordToken, database)

	server := api.NewServer(cfg.Port, discordBot)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Web server failed: %v", err)
		}
	}()

	ctx := context.Background()
	if err := discordBot.Start(ctx); err != nil {
		log.Fatalf("Failed to start Discord bot: %v", err)
	}

	log.Println("PulseKeep (Go version) is running with sharding. Press CTRL+C to exit.")

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s

	log.Println("Shutdown signal received. Initiating graceful cleanup...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down web server: %v", err)
	} else {
		log.Println("Web server shut down successfully.")
	}

	log.Println("Stopping Discord bot connection...")
	discordBot.Stop(shutdownCtx)
	log.Println("Discord bot connection stopped cleanly. Farewell!")
}
