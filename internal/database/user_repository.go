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
	query := `INSERT INTO user (userName, email, password, createdAt) VALUES (?, ?, ?, ?)`

	result, err := ExecWithTimeout(query, username, email, password, time.Now())
	if err != nil {
		log.Printf("SQL error when creating user: %v", err)
		return fmt.Errorf("failed to create user: %v", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("no rows affected, user not created")
	}
	return nil
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	var createdAtStr string

	query := `SELECT * FROM user WHERE email = ?`
	err := GetDB().QueryRow(query, email).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&createdAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user with email %s not found", email)
		}
		return models.User{}, fmt.Errorf("database error when fetching user: %v", err)
	}

	if createdAtStr != "" {
		if parsed, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			user.CreatedAt = parsed
		}
	}

	return user, nil
}
