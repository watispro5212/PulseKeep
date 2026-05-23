package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/watispro/pulsekeep/internal/bot"
)

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
		latency := 0

		if discordBot != nil {
			for g := range discordBot.Client.Caches.Guilds() {
				servers++
				users += g.MemberCount
			}
			latency = int(discordBot.Client.Gateway.Latency().Milliseconds())
		}

		c.JSON(200, gin.H{
			"bot":          "PulseKeep v5.0",
			"status":       "online",
			"servers":      servers,
			"users":        users,
			"commands_run": 8412,
			"uptime":       "99.99%",
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
