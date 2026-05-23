package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/watispro/pulsekeep/internal/db"
)

type Server struct {
	httpServer *http.Server
	port       string
	startedAt  time.Time
	database   *db.Database
}

func NewServer(port string, allowedOrigin string, database *db.Database) *Server {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	startedAt := time.Now()
	r := gin.Default()

	// CORS Middleware to allow Netlify frontend to securely fetch stats from Fly.io backend
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
		})
	})

	return &Server{
		httpServer: &http.Server{
			Addr:              ":" + port,
			Handler:           r,
			ReadHeaderTimeout: 5 * time.Second,
		},
		port:      port,
		startedAt: startedAt,
		database:  database,
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
