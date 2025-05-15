package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

// Setup the database schema for the jobs table
func SetupDatabase(conn *pgx.Conn) error {
	query := `
	CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		repo_url TEXT NOT NULL,
		branch TEXT,
		status TEXT,
		created_at TIMESTAMP DEFAULT now(),
		finished_at TIMESTAMP,
		duration_ms INTEGER,
		lint_errors INTEGER,
		lint_warnings INTEGER,
		build_log TEXT
	);
	`
	_, err := conn.Exec(context.Background(), query)
	if err != nil {
		log.Printf("Failed to setup jobs table: %+v", err)
		return err
	}
	log.Printf("Jobs table setup!")
	return nil
}

// Connect to the database
func Connect() (*pgx.Conn, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
