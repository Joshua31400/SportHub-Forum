package database

import (
	"SportHub-Forum/internal/models"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// CreateUser inserts a new user into the database
func CreateUser(username, email, password string) error {
	// Current time for the createdAt field
	createdAt := time.Now()

	// Query to insert a new user
	query := `INSERT INTO user (userName, email, password, createdAt) VALUES (?, ?, ?, ?)`

	// Using timeout function instead of direct db.Exec
	result, err := ExecWithTimeout(query, username, email, password, createdAt)
	if err != nil {
		log.Printf("SQL error when creating user: %v", err)
		return fmt.Errorf("failed to create user: %v", err)
	}

	// Check if the insertion was successful
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, user not created")
	}
	return nil
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := `SELECT userID, userName, email, password, createdAt FROM user WHERE email = ?`

	// Using timeout function instead of direct db.QueryRow
	err := QueryRowWithTimeout(query, email).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user with email %s not found", email)
		}
		return models.User{}, fmt.Errorf("database error when fetching user: %v", err)
	}

	return user, nil
}
