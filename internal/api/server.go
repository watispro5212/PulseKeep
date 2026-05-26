package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/watispro5212/PulseKeep/internal/cache"
	"github.com/watispro5212/PulseKeep/internal/db"
)

type Server struct {
	httpServer *http.Server
	port       string
	startedAt  time.Time
	database   *db.Database
	cache      *cache.Cache
}

func NewServer(port string, allowedOrigin string, database *db.Database, memCache *cache.Cache) *Server {
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

		var goVersion string
		if memCache != nil {
			goVersion = runtime.Version()
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      status,
			"database":    dbStatus,
			"uptime":      formatDuration(time.Since(startedAt)),
			"bot_uptime":  formatUptime(memCache),
			"go_version":  goVersion,
			"servers":     getGuildCount(memCache),
			"users":       getUserCount(memCache),
			"commands":    getCommandsRun(memCache),
		})
	})

	r.GET("/stats", func(c *gin.Context) {
		servers := getGuildCount(memCache)
		users := getUserCount(memCache)
		cmds := getCommandsRun(memCache)

		c.JSON(http.StatusOK, gin.H{
			"bot":          "PulseKeep v5.8",
			"status":       "online",
			"servers":      servers,
			"users":        users,
			"commands_run": cmds,
			"uptime":       formatDuration(time.Since(startedAt)),
			"go_version":   runtime.Version(),
			"features": []string{
				"moderation",
				"tickets",
				"audit_logs",
				"economy",
				"shop",
				"polls",
				"gambling",
				"fishing",
				"mining",
				"channel_management",
			},
		})
	})

	r.GET("/api/dashboard", func(c *gin.Context) {
		servers := getGuildCount(memCache)
		users := getUserCount(memCache)
		cmds := getCommandsRun(memCache)

		economyStats := gin.H{"status": "unavailable"}
		if database != nil && database.Conn != nil {
			var totalRecords, totalBalance, avgBalance int
			var topUsers []gin.H

			err := database.Conn.QueryRow(`SELECT COUNT(*), COALESCE(SUM(balance),0), COALESCE(AVG(balance),0) FROM user_economy`).Scan(&totalRecords, &totalBalance, &avgBalance)
			if err == nil {
				rows, err := database.Conn.Query(`SELECT user_id, name, balance FROM user_economy ORDER BY balance DESC LIMIT 5`)
				if err == nil {
					for rows.Next() {
						var uid, name string
						var bal int
						if rows.Scan(&uid, &name, &bal) == nil {
							topUsers = append(topUsers, gin.H{"id": uid, "name": name, "balance": bal})
						}
					}
					rows.Close()
				}
				economyStats = gin.H{
					"status":        "ok",
					"total_records": totalRecords,
					"total_balance": totalBalance,
					"avg_balance":   avgBalance,
					"top_users":     topUsers,
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"bot":          "PulseKeep v5.8",
			"status":       "online",
			"servers":      servers,
			"users":        users,
			"commands_run": cmds,
			"uptime":       formatDuration(time.Since(startedAt)),
			"bot_uptime":   formatUptime(memCache),
			"go_version":   runtime.Version(),
			"database":     economyStats,
			"features": []string{
				"moderation",
				"tickets",
				"audit_logs",
				"economy",
				"shop",
				"polls",
				"gambling",
				"fishing",
				"mining",
				"channel_management",
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
		cache:     memCache,
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
	secs := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, secs)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, secs)
	}
	return fmt.Sprintf("%ds", secs)
}

func formatUnit(value int, suffix string) string {
	if value <= 0 {
		return ""
	}
	return strings.TrimSpace(strings.Join([]string{strconv.Itoa(value), suffix}, ""))
}

func getGuildCount(memCache *cache.Cache) int64 {
	if memCache == nil {
		return 0
	}
	return memCache.GuildCount.Load()
}

func getUserCount(memCache *cache.Cache) int64 {
	if memCache == nil {
		return 0
	}
	return memCache.UserCount.Load()
}

func getCommandsRun(memCache *cache.Cache) int64 {
	if memCache == nil {
		return 0
	}
	return memCache.CommandsRun.Load()
}

func formatUptime(memCache *cache.Cache) string {
	if memCache == nil || memCache.StartedAt.IsZero() {
		return "unknown"
	}
	return formatDuration(time.Since(memCache.StartedAt))
}
