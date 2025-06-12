package database

import (
	"SportHub-Forum/internal/authentification"
	"SportHub-Forum/internal/models"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// CreateUser inserts a new user into the database with hashed password
func CreateUser(username, email, password string) error {
	hashedPassword, err := authentification.HashPassword(password)
	if err != nil {
		log.Printf("Ash Error: %v", err)
		return fmt.Errorf("Wrong password: %v", err)
	}

	query := `INSERT INTO user (userName, email, password, createdAt) VALUES (?, ?, ?, ?)`

	// Use ashed password and current time for createdAt
	result, err := ExecWithTimeout(query, username, email, hashedPassword, time.Now())
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

func GetUserByID(db *sql.DB, userID int) (models.User, error) {
	var user models.User
	var createdAtStr string

	query := `SELECT userID, userName, email, password, createdAt FROM user WHERE userID = ?`
	err := db.QueryRow(query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&createdAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("User ID: %d not found", userID)
		}
		return models.User{}, fmt.Errorf("Error to get user: %v", err)
	}

	if createdAtStr != "" {
		formats := []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02T15:04:05Z"}
		for _, format := range formats {
			if parsed, err := time.Parse(format, createdAtStr); err == nil {
				user.CreatedAt = parsed
				break
			}
		}
	}
	return user, nil
}
