package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err = connectWithRetry(); err != nil {
		return err
	}

	go startPeriodicPing()

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

// QueryRowWithTimeout executes a QueryRow with a 5-second timeout
func QueryRowWithTimeout(query string, args ...interface{}) *sql.Row {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	return db.QueryRowContext(ctx, query, args...)
}

// QueryWithTimeout executes a Query with a 5-second timeout
func QueryWithTimeout(query string, args ...interface{}) (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	return db.QueryContext(ctx, query, args...)
}

// connectWithRetry attempts to connect to the database with retries
func connectWithRetry() error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = db.Ping()
		if err == nil {
			return nil
		}
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
			if err := connectWithRetry(); err == nil {
				// Reconnected successfully
			}
		}
	}
}

// InitTables verifies and initializes database tables if needed
func InitTables() error {
	var exists int
	err := db.QueryRow("SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'users' LIMIT 1").Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error verifying if table exists: %v", err)
	}

	return nil
}
