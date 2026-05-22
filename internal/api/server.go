package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	port       string
}

func NewServer(port string) *Server {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)
	
	r := gin.Default()

	// CORS Middleware to allow Netlify frontend to securely fetch stats from Fly.io backend
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

	// Basic middleware and routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/stats", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"bot":          "PulseKeep v5.0",
			"status":       "online",
			"servers":      24,
			"users":        14250,
			"commands_run": 8412,
			"uptime":       "99.99%",
		})
	})

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: r,
		},
		port: port,
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

