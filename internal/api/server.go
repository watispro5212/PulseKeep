package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/watispro5212/PulseKeep/internal/auth"
	"github.com/watispro5212/PulseKeep/internal/bot/automod"
	"github.com/watispro5212/PulseKeep/internal/cache"
	"github.com/watispro5212/PulseKeep/internal/config"
	"github.com/watispro5212/PulseKeep/internal/db"
)

type Server struct {
	httpServer *http.Server
	config     *config.Config
	database   *db.Database
	cache      *cache.Cache
	cfgStore   *automod.ConfigStore
}

const discordAPIURL = "https://discord.com/api/v10"

func NewServer(cfg *config.Config, database *db.Database, memCache *cache.Cache, cfgStore *automod.ConfigStore) *Server {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	startedAt := time.Now()
	r := gin.Default()

	// CORS Middleware to allow Netlify frontend to securely fetch stats from Fly.io backend
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if cfg.AllowedOrigin == "*" || origin == cfg.AllowedOrigin || isAllowedOrigin(origin, cfg.AllowedOrigin) {
			if cfg.AllowedOrigin == "*" {
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

		goVersion := runtime.Version()

		httpStatus := http.StatusOK
		if status == "degraded" {
			httpStatus = http.StatusServiceUnavailable
		}
		c.JSON(httpStatus, gin.H{
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
"bot": "PulseKeep v6.0.0",
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
			var totalRecords, totalBalance int
			var avgBalance float64
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
"bot": "PulseKeep v6.0.0",
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

	// Discord OAuth2 endpoints
	r.GET("/auth/discord", func(c *gin.Context) {
		if cfg.DiscordClientID == "" || cfg.DiscordClientSecret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "OAuth not configured"})
			return
		}
		// In a real app, you'd generate and store a state value for CSRF protection
		params := url.Values{}
		params.Set("client_id", cfg.DiscordClientID)
		params.Set("redirect_uri", cfg.DiscordRedirectURI)
		params.Set("response_type", "code")
		params.Set("scope", "identify guilds")
		redirectURL := fmt.Sprintf("%s/oauth2/authorize?%s", discordAPIURL, params.Encode())
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	})

	r.GET("/auth/discord/callback", func(c *gin.Context) {
		code := c.Query("code")
		// In a real app, you'd verify state matches the one stored in session/cookie
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code parameter"})
			return
		}

		token, err := auth.ExchangeCode(cfg, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token: " + err.Error()})
			return
		}

		user, err := auth.GetUser(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info: " + err.Error()})
			return
		}

		// Redirect to dashboard with token in URL fragment
		dashURL := fmt.Sprintf("/dashboard.html#access_token=%s&user_id=%s", url.QueryEscape(token), url.QueryEscape(user.ID))
		c.Redirect(http.StatusTemporaryRedirect, dashURL)
	})

	// Guild config endpoints
	r.GET("/api/guilds", func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token parameter"})
			return
		}
		guilds, err := auth.GetUserGuilds(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch guilds: " + err.Error()})
			return
		}
		user, err := auth.GetUser(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user, "guilds": guilds})
	})

	r.GET("/api/guild/:id/config", func(c *gin.Context) {
		guildID := c.Param("id")
		// TODO: Verify user has permission to manage this guild (via Discord API)
		if cfgStore == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Config store not initialized"})
			return
		}
		guildCfg := cfgStore.Get(guildID)
		c.JSON(http.StatusOK, guildCfg)
	})

	r.POST("/api/guild/:id/config", func(c *gin.Context) {
		guildID := c.Param("id")
		// TODO: Verify user has permission to manage this guild (via Discord API and OAuth2 token)
		var configUpdate automod.GuildConfig
		if err := c.ShouldBindJSON(&configUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}
		// Ensure GuildID matches
		configUpdate.GuildID = guildID
		if cfgStore == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Config store not initialized"})
			return
		}
		if err := cfgStore.Update(&configUpdate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update config: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return &Server{
		httpServer: &http.Server{
			Addr:              ":" + cfg.Port,
			Handler:           r,
			ReadHeaderTimeout: 5 * time.Second,
		},
		config:   cfg,
		database: database,
		cache:    memCache,
		cfgStore: cfgStore,
	}
}

func (s *Server) Start() error {
	log.Printf("Starting web server on port %s", s.config.Port)
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
