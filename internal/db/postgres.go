package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	Conn *sql.DB
}

func Connect(databaseURL string) *Database {
	conn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database!")
	return &Database{Conn: conn}
}

// Ensure the connection is closed when the application shuts down
func (db *Database) Close() error {
	return db.Conn.Close()
}
