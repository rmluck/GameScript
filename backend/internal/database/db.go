// Internal database package for managing PostgreSQL connections

package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)


type DB struct {
	Conn *sql.DB
}

func NewConnection() (*DB, error) {
	// Try to get database url first (for production)
	databaseURL := os.Getenv("DATABASE_URL")

	// If database URL not set, build database configuration from individual environment variables (for local development)
	if databaseURL == "" {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbName := os.Getenv("DB_NAME")
		dbSSLMode := os.Getenv("DB_SSL_MODE")

		// Construct the connection string
		databaseURL = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbName, dbSSLMode)
	}
	
	
	// Log the DSN for debugging purposes (avoid logging sensitive info in production)
	log.Printf("Connecting to database...")
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Verify the connection is valid
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	log.Println("Database connected successfully.")

	return &DB{Conn: db}, nil
}

func (db *DB) Close() error {
	return db.Conn.Close()
}

// Executes a query and returns the resulting rows
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Conn.Query(query, args...)
}