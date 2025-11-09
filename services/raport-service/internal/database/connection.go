package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	OrdersDB  *sql.DB
	RaportsDB *sql.DB
}

func NewConnection() (*DB, error) {
	ordersDB, err := connectToOrdersDB()
	if err != nil {
		return nil, fmt.Errorf("ERROR: nie można połączyć się z bazą danych: %v", err)
	}
	reportsDB, err := connectToRaportsDB()
	if err != nil {
		return nil, fmt.Errorf("ERROR: nie można połączyć się z bazą danych raportów: %v", err)
	}

	log.Println("Pomyślnie połączono z bazą danych")
	return &DB{OrdersDB: ordersDB, RaportsDB: reportsDB}, nil
}

func connectToOrdersDB() (*sql.DB, error) {
	host := getEnv("ORDERS_DB_HOST", "localhost")
	port := getEnv("ORDERS_DB_PORT", "5432")
	user := getEnv("ORDERS_DB_USER", "postgres")
	password := getEnv("ORDERS_DB_PASSWORD", "password123")
	dbname := getEnv("ORDERS_DB_NAME", "orders_management")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToRaportsDB() (*sql.DB, error) {
	host := getEnv("RAPORTS_DB_HOST", "localhost")
	port := getEnv("RAPORTS_DB_PORT", "5432")
	user := getEnv("RAPORTS_DB_USER", "postgres")
	password := getEnv("RAPORTS_DB_PASSWORD", "password123")
	dbname := getEnv("RAPORTS_DB_NAME", "raports_management")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func (db *DB) Close() {
	if db.OrdersDB != nil {
		db.OrdersDB.Close()
	}
	if db.RaportsDB != nil {
		db.RaportsDB.Close()
	}
}
