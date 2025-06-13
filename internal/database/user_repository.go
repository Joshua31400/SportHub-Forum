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
		log.Printf("Hash Error: %v", err)
		return fmt.Errorf("Wrong password: %v", err)
	}

	query := `INSERT INTO user (userName, email, password, auth_provider, is_verified, createdAt, updated_at) 
              VALUES (?, ?, ?, 'local', FALSE, ?, ?)`

	currentTime := time.Now()
	result, err := ExecWithTimeout(query, username, email, hashedPassword, currentTime, currentTime)
	if err != nil {
		log.Printf("SQL error when creating user: %v", err)
		return fmt.Errorf("failed to create user: %v", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("no rows affected, user not created")
	}
	return nil
}

// CreateGitHubUser creates a new user from GitHub OAuth
func CreateGitHubUser(email, username, githubID, avatar string) (*models.User, error) {
	query := `INSERT INTO user (userName, email, github_id, avatar, auth_provider, is_verified, createdAt, updated_at)
              VALUES (?, ?, ?, ?, 'github', TRUE, ?, ?)`

	currentTime := time.Now()
	result, err := GetDB().Exec(query, username, email, githubID, avatar, currentTime, currentTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub user: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %v", err)
	}

	return &models.User{
		UserID:       int(userID),
		Username:     username,
		Email:        email,
		GitHubID:     githubID,
		Avatar:       avatar,
		AuthProvider: "github",
		IsVerified:   true,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}, nil
}

// GetUserByGitHubID retrieves a user by their GitHub ID
func GetUserByGitHubID(githubID string) (*models.User, error) {
	var user models.User
	var createdAtStr, updatedAtStr string
	var password, avatar sql.NullString
	var authProvider sql.NullString
	var isVerified sql.NullBool
	var githubIDDB sql.NullString

	query := `SELECT userID, userName, email, password, 
                 DATE_FORMAT(createdAt, '%Y-%m-%d %H:%i:%s') as createdAt,
                 github_id, avatar, auth_provider, is_verified,
                 DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') as updated_at
              FROM user WHERE github_id = ?`

	err := GetDB().QueryRow(query, githubID).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&password,
		&createdAtStr,
		&githubIDDB,
		&avatar,
		&authProvider,
		&isVerified,
		&updatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving user by GitHub ID: %v", err)
	}

	// Parse dates and handle nulls
	if createdAtStr != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			user.CreatedAt = parsed
		}
	}

	if password.Valid {
		user.Password = password.String
	}

	if githubIDDB.Valid {
		user.GitHubID = githubIDDB.String
	}

	if avatar.Valid {
		user.Avatar = avatar.String
	}

	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "local"
	}

	if isVerified.Valid {
		user.IsVerified = isVerified.Bool
	}

	if updatedAtStr != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
			user.UpdatedAt = parsed
		}
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by their email (mise à jour pour gérer les nouvelles colonnes)
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	var createdAtStr, updatedAtStr string
	var password, avatar sql.NullString
	var authProvider sql.NullString
	var isVerified sql.NullBool
	var githubID sql.NullString

	query := `SELECT userID, userName, email, password, 
                 DATE_FORMAT(createdAt, '%Y-%m-%d %H:%i:%s') as createdAt,
                 github_id, avatar, auth_provider, is_verified,
                 DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') as updated_at
              FROM user WHERE email = ?`

	err := GetDB().QueryRow(query, email).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&password,
		&createdAtStr,
		&githubID,
		&avatar,
		&authProvider,
		&isVerified,
		&updatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user with email %s not found", email)
		}
		return models.User{}, fmt.Errorf("database error when fetching user: %v", err)
	}

	// Parse dates and handle nulls (même logique que GetUserByGitHubID)
	if createdAtStr != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			user.CreatedAt = parsed
		}
	}

	if password.Valid {
		user.Password = password.String
	}

	if githubID.Valid {
		user.GitHubID = githubID.String
	}

	if avatar.Valid {
		user.Avatar = avatar.String
	}

	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "local"
	}

	if isVerified.Valid {
		user.IsVerified = isVerified.Bool
	}

	if updatedAtStr != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
			user.UpdatedAt = parsed
		}
	}

	return user, nil
}

// GetUserByID mise à jour pour gérer les nouvelles colonnes
func GetUserByID(db *sql.DB, userID int) (models.User, error) {
	var user models.User
	var createdAtStr, updatedAtStr string
	var password, avatar sql.NullString
	var authProvider sql.NullString
	var isVerified sql.NullBool
	var githubID sql.NullString

	query := `SELECT userID, userName, email, password, 
                 DATE_FORMAT(createdAt, '%Y-%m-%d %H:%i:%s') as createdAt,
                 github_id, avatar, auth_provider, is_verified,
                 DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') as updated_at
              FROM user WHERE userID = ?`

	err := db.QueryRow(query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&password,
		&createdAtStr,
		&githubID,
		&avatar,
		&authProvider,
		&isVerified,
		&updatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("User ID: %d not found", userID)
		}
		return models.User{}, fmt.Errorf("Error to get user: %v", err)
	}

	// Parse dates and handle nulls (même logique)
	if createdAtStr != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			user.CreatedAt = parsed
		}
	}

	if password.Valid {
		user.Password = password.String
	}

	if githubID.Valid {
		user.GitHubID = githubID.String
	}

	if avatar.Valid {
		user.Avatar = avatar.String
	}

	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "local"
	}

	if isVerified.Valid {
		user.IsVerified = isVerified.Bool
	}

	if updatedAtStr != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
			user.UpdatedAt = parsed
		}
	}

	return user, nil
}
