package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/watispro5212/PulseKeep/internal/api"
	"github.com/watispro5212/PulseKeep/internal/bot"
	"github.com/watispro5212/PulseKeep/internal/bot/automod"
	"github.com/watispro5212/PulseKeep/internal/cache"
	"github.com/watispro5212/PulseKeep/internal/config"
	"github.com/watispro5212/PulseKeep/internal/db"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Init Database
	database := db.Connect(cfg.DatabaseURL)
	defer database.Close()

	// 3. Init Cache
	memCache := cache.New()

	// 4. Init Bot (to get config store)
	var discordBot *bot.Bot

	if cfg.BotDisabled {
		log.Println("Discord bot is disabled; API endpoints are still available.")
	} else {
		var dbConn *sql.DB
		if database != nil {
			dbConn = database.Conn
		}
		discordBot = bot.New(cfg.DiscordToken, memCache, dbConn)
	}

	// 5. Init Web Server
	var cfgStore *automod.ConfigStore
	if discordBot != nil {
		cfgStore = discordBot.GetConfigStore()
	}
	server := api.NewServer(cfg, database, memCache, cfgStore)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Web server failed: %v", err)
		}
	}()

	// 6. Start Bot
	if !cfg.BotDisabled && discordBot != nil {
		ctx := context.Background()
		if err := discordBot.Start(ctx); err != nil {
			log.Fatalf("Failed to start Discord bot: %v", err)
		}
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
	if discordBot != nil {
		log.Println("Stopping Discord bot connection...")
		discordBot.Stop(shutdownCtx)
		log.Println("Discord bot connection stopped cleanly.")
	}
	log.Println("PulseKeep stopped cleanly. Farewell!")
}
