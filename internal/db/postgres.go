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
		return nil
	}

	conn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(30 * time.Minute)

	if err := conn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database!")
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
