package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofrs/uuid/v5"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

func main() {
	// Parse command-line flags
	config := Config{}
	flag.StringVar(&config.Host, "host", "localhost", "PostgreSQL host")
	flag.IntVar(&config.Port, "port", 5432, "PostgreSQL port")
	flag.StringVar(&config.User, "user", "postgres", "PostgreSQL user")
	flag.StringVar(&config.Password, "password", "", "PostgreSQL password")
	flag.StringVar(&config.Database, "database", "postgres", "PostgreSQL database name")
	flag.StringVar(&config.SSLMode, "sslmode", "disable", "PostgreSQL SSL mode (disable, require, verify-ca, verify-full)")

	flag.Parse()

	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode,
	)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Printf("Successfully connected to PostgreSQL at %s:%d/%s", config.Host, config.Port, config.Database)

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start inserting
	log.Println("Starting continuous inserts... Press Ctrl+C to stop")

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			log.Println("Received shutdown signal, exiting gracefully...")
			return
		case <-ticker.C:
			if err := insertRecord(db); err != nil {
				log.Printf("Failed to insert record: %v", err)
			}
		}
	}
}

func insertRecord(db *sql.DB) error {
	// Generate UUID v7 (time-ordered UUID)
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate UUID v7: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO events (id) VALUES ($1)`
	_, err = db.ExecContext(ctx, query, id)

	return err
}
