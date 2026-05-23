package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
<<<<<<< HEAD
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/watispro/pulsekeep/internal/db"
=======
	"time"

	"github.com/gin-gonic/gin"
	"github.com/watispro/pulsekeep/internal/bot"
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
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
<<<<<<< HEAD
	startedAt  time.Time
	database   *db.Database
}

func NewServer(port string, allowedOrigin string, database *db.Database) *Server {
	// Set Gin to release mode in production
=======
	Bot        *bot.Bot
}

func NewServer(port string, discordBot *bot.Bot) *Server {
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
	gin.SetMode(gin.ReleaseMode)

	startedAt := time.Now()
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowedOrigin == "*" || origin == allowedOrigin || isAllowedOrigin(origin, allowedOrigin) {
			if allowedOrigin == "*" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Vary", "Origin")
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		status := "ok"
		dbStatus := "not_configured"
		if database != nil {
			dbStatus = "ok"
			if err := database.Ping(); err != nil {
				status = "degraded"
				dbStatus = "unavailable"
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   status,
			"database": dbStatus,
			"uptime":   formatDuration(time.Since(startedAt)),
		})
	})

	r.GET("/stats", func(c *gin.Context) {
<<<<<<< HEAD
		c.JSON(http.StatusOK, gin.H{
			"bot":          "PulseKeep v5.0",
			"status":       "online",
			"servers":      24,
			"users":        14250,
			"commands_run": 8412,
			"uptime":       formatDuration(time.Since(startedAt)),
			"latency_ms":   9,
			"features": []string{
				"moderation",
				"tickets",
				"audit_logs",
				"economy",
			},
=======
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
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
		})
	})

	return &Server{
		httpServer: &http.Server{
			Addr:              ":" + port,
			Handler:           r,
			ReadHeaderTimeout: 5 * time.Second,
		},
<<<<<<< HEAD
		port:      port,
		startedAt: startedAt,
		database:  database,
=======
		port: port,
		Bot:  discordBot,
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
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
<<<<<<< HEAD

func isAllowedOrigin(origin string, allowedOrigins string) bool {
	if origin == "" || allowedOrigins == "" {
		return false
	}

	for _, allowed := range strings.Split(allowedOrigins, ",") {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}
	return false
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return strings.TrimSpace(strings.Join([]string{
			formatUnit(days, "d"),
			formatUnit(hours, "h"),
			formatUnit(minutes, "m"),
		}, " "))
	}
	if hours > 0 {
		return strings.TrimSpace(strings.Join([]string{
			formatUnit(hours, "h"),
			formatUnit(minutes, "m"),
		}, " "))
	}
	if minutes <= 0 {
		return "0m"
	}
	return formatUnit(minutes, "m")
}

func formatUnit(value int, suffix string) string {
	if value <= 0 {
		return ""
	}
	return strings.TrimSpace(strings.Join([]string{strconv.Itoa(value), suffix}, ""))
}
=======
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
