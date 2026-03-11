package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaSQL string

// DB wraps the database connection
type DB struct {
	*sql.DB
}

// New creates a new database connection and initializes the schema
func New(dataDir string) (*DB, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "app.db")
	
	// Open database connection
	sqlDB, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{sqlDB}

	// Initialize schema
	if err := db.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Printf("Database initialized at %s", dbPath)
	return db, nil
}

// initSchema creates all tables if they don't exist
func (db *DB) initSchema() error {
	_, err := db.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// BeginTx starts a new transaction
func (db *DB) BeginTx() (*sql.Tx, error) {
	return db.Begin()
}

