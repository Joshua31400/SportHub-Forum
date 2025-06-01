package database

import (
	"SportHub-Forum/internal/models"
	"fmt"
	"log"
	"time"
)

func CreateUser(username, email, password string) error {
	if err := db.Ping(); err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return fmt.Errorf("Database unrichable: %v", err)
	}

	// Get the current time for the createdAt field
	createdAt := time.Now()

	// Query to insert a new user
	query := `INSERT INTO user (userName, email, password, createdAt) VALUES (?, ?, ?, ?)`

	result, err := db.Exec(query, username, email, password, createdAt)
	if err != nil {
		log.Printf("Error SQL for create user: %v", err)
		return fmt.Errorf("Error to create user: %v", err)
	}

	// Chef if the insertion was successful
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected, user not created")
	}
	return nil
}

// Get user by email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := `SELECT userID, userName, email, password, createdAt FROM user WHERE email = ?`

	err := db.QueryRow(query, email).Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("User not found: %v", err)
	}

	return user, nil
}
