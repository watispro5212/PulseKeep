package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/watispro/pulsekeep/internal/bot"
)

// formatDuration returns a human-readable uptime string
func formatDuration(d time.Duration) string {
	total := int(d.Seconds())
	days := total / 86400
	hours := (total % 86400) / 3600
	minutes := (total % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

type Server struct {
	httpServer *http.Server
	port       string
	Bot        *bot.Bot
}

func NewServer(port string, discordBot *bot.Bot) *Server {
	gin.SetMode(gin.ReleaseMode)
	
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/stats", func(c *gin.Context) {
		servers := 0
		users := 0
		commandsRun := 0
		latency := 0
		uptime := "—"

		if discordBot != nil {
			// Count servers and users from the cache
			for g := range discordBot.Client.Caches.Guilds() {
				servers++
				users += g.MemberCount
			}
			// Get gateway latency
			latency = int(discordBot.Client.Gateway.Latency().Milliseconds())
			// Calculate uptime from bot start time
			startTime := discordBot.StartTime
			if !startTime.IsZero() {
				duration := time.Since(startTime)
				uptime = formatDuration(duration)
			}
			// Get command run count from database if available
			if discordBot.DB != nil {
				count, err := discordBot.DB.GetCommandsRunCount(c.Request.Context())
				if err == nil {
					commandsRun = count
				}
			}
		}

		c.JSON(200, gin.H{
			"bot":          "PulseKeep v5.0",
			"status":       "online",
			"servers":      servers,
			"users":        users,
			"commands_run": commandsRun,
			"uptime":       uptime,
			"latency":      latency,
		})
	})

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: r,
		},
		port: port,
		Bot:  discordBot,
	}
}

func (s *Server) Start() error {
	log.Printf("Starting web server on port %s", s.port)
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down web server gracefully...")
	return s.httpServer.Shutdown(ctx)
}
