package db

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	Conn *sql.DB
}

func Connect(databaseURL string) *Database {
	if databaseURL == "" {
		log.Println("DATABASE_URL is not set; running without database.")
		return nil
	}

	conn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(30 * time.Minute)

	if err := conn.Ping(); err != nil {
		log.Printf("Failed to ping database (will retry later): %v", err)
		// Return the connection anyway; health check will report degraded
		// and subsequent operations will reconnect automatically
	}

	log.Println("Database connection established.")
	database := &Database{Conn: conn}
	database.Migrate()
	return database
}

// Ensure the connection is closed when the application shuts down
func (db *Database) Close() error {
	if db == nil || db.Conn == nil {
		return nil
	}
	return db.Conn.Close()
}

func (db *Database) Ping() error {
	if db == nil || db.Conn == nil {
		return errors.New("database not configured")
	}
	return db.Conn.Ping()
}
