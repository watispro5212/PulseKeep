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
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Init Database
	database := db.Connect(cfg.DatabaseURL)
	defer database.Close()

	// 3. Init Cache
	memCache := cache.New()
	_ = memCache // We will pass this to handlers later

	// 4. Init Web Server
	server := api.NewServer(cfg.Port)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Web server failed: %v", err)
		}
	}()

	// 5. Init Bot
	discordBot := bot.New(cfg.DiscordToken)

	ctx := context.Background()
	if err := discordBot.Start(ctx); err != nil {
		log.Fatalf("Failed to start Discord bot: %v", err)
	}

	log.Println("PulseKeep (Go version) is running. Press CTRL+C to exit.")

	// Wait for shutdown signal (SIGINT / SIGTERM / os.Interrupt)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s

	log.Println("Shutdown signal received. Initiating graceful cleanup...")

	// Create a context with 5s timeout to allow ongoing operations to drain cleanly
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Gracefully shut down the web server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down web server: %v", err)
	} else {
		log.Println("Web server shut down successfully.")
	}

	// 2. Stop the Discord Bot connection
	log.Println("Stopping Discord bot connection...")
	discordBot.Stop(shutdownCtx)
	log.Println("Discord bot connection stopped cleanly. Farewell!")
}
