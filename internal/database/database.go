package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var db *sql.DB

const (
	maxRetries   = 10
	retryDelay   = time.Second * 3
	queryTimeout = time.Second * 5
)

// InitDB initializes the database connection with proper configuration
func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %v", err)
	}

	// Pool settings
	db.SetMaxOpenConns(20)                 // Limit to 20 simultaneous connections
	db.SetMaxIdleConns(10)                 // Keep up to 10 idle connections
	db.SetConnMaxLifetime(time.Minute * 5) // Maximum lifetime of a connection

	if err = connectWithRetry(); err != nil {
		return err
	}

	go startPeriodicPing()

	log.Println("Connected to the database successfully")
	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// ExecWithTimeout executes a query with a 5-second timeout
func ExecWithTimeout(query string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	return db.ExecContext(ctx, query, args...)
}

// connectWithRetry attempts to connect to the database with retries
func connectWithRetry() error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = db.Ping()
		if err == nil {
			return nil // Connection successful
		}

		log.Printf("Connection attempt failed (%d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("connection failed after %d attempts: %v", maxRetries, err)
}

// startPeriodicPing checks database connection every 30 seconds
func startPeriodicPing() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := db.Ping(); err != nil {
			log.Printf("Periodic ping failed: %v. Attempting reconnection...", err)
			if err := connectWithRetry(); err != nil {
				log.Printf("Reconnection failed: %v", err)
			} else {
				log.Println("Successfully reconnected to database")
			}
		}
	}
}
