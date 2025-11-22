package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewConnection() (*DB, error) {
	databaseURL := os.Getenv("AUTH_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:password@localhost:5433/auth_service?sslmode=disable"
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Testowanie połączenia
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("błąd połączenia z bazą danych: %w", err)
	}

	log.Println("Pomyślnie połączono z bazą danych")
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// Tworzenie tabeli users
func (db *DB) CreateUsersTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL DEFAULT 'customer',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Create index for faster email lookups
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("błąd tworzenia tabeli users: %w", err)
	}

	log.Println("Tabela users utworzona pomyślnie")
	return nil
}
