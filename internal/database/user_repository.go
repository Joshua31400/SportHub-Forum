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

	query := `INSERT INTO user (userName, email, password, auth_provider, createdAt) VALUES (?, ?, ?, 'local', ?)`

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
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	var createdAtStr string
	var updatedAtStr sql.NullString
	var password sql.NullString
	var avatar sql.NullString
	var authProvider sql.NullString
	var isVerified sql.NullBool
	var googleID sql.NullString

	query := `SELECT 
        userid, username, email, password, 
        DATE_FORMAT(createdat, '%Y-%m-%d %H:%i:%s') as createdat,
        google_id, avatar, auth_provider, is_verified,
        DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') as updated_at
    FROM user WHERE email = ?`

	err := GetDB().QueryRow(query, email).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&password,
		&createdAtStr,
		&googleID,
		&avatar,
		&authProvider,
		&isVerified,
		&updatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("database error when fetching user: %v", err)
	}

	if createdAtStr != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			user.CreatedAt = parsed
		}
	}

	if updatedAtStr.Valid && updatedAtStr.String != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", updatedAtStr.String); err == nil {
			user.UpdatedAt = parsed
		}
	}

	if password.Valid {
		user.Password = password.String
	}

	if googleID.Valid {
		user.GoogleID = googleID.String
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

	return &user, nil
}

// GetUserByGitHubID retrieves a user by their GitHub ID
func GetUserByGitHubID(githubID string) (*models.User, error) {
	var user models.User
	var createdAtStr string
	var password sql.NullString
	var avatar sql.NullString
	var authProvider sql.NullString
	var isVerified sql.NullBool
	var updatedAtStr sql.NullString
	var googleID sql.NullString
	var githubIDDB sql.NullString

	query := `SELECT userID, username, email, password, 
                 DATE_FORMAT(createdAt, '%Y-%m-%d %H:%i:%s') as createdAt,
                 google_id, github_id, avatar, auth_provider, is_verified,
                 DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') as updated_at
          FROM user WHERE github_id = ?`

	err := GetDB().QueryRow(query, githubID).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&password,
		&createdAtStr,
		&googleID,
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

	if createdAtStr != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			user.CreatedAt = parsed
		}
	}

	if password.Valid {
		user.Password = password.String
	} else {
		user.Password = ""
	}

	if googleID.Valid {
		user.GoogleID = googleID.String
	} else {
		user.GoogleID = ""
	}

	if githubIDDB.Valid {
		user.GitHubID = githubIDDB.String
	} else {
		user.GitHubID = ""
	}

	if avatar.Valid {
		user.Avatar = avatar.String
	} else {
		user.Avatar = ""
	}

	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "local"
	}

	if isVerified.Valid {
		user.IsVerified = isVerified.Bool
	} else {
		user.IsVerified = false
	}

	if updatedAtStr.Valid && updatedAtStr.String != "" {
		if parsed, err := time.Parse("2006-01-02 15:04:05", updatedAtStr.String); err == nil {
			user.UpdatedAt = parsed
		}
	}

	return &user, nil
}

// CreateGitHubUser creates a new user from GitHub OAuth
func CreateGitHubUser(email, username, githubID, avatar string) (*models.User, error) {
	query := `INSERT INTO user (username, email, github_id, avatar, auth_provider, is_verified, createdAt, updated_at)
              VALUES (?, ?, ?, ?, 'github', TRUE, NOW(), NOW())`

	result, err := GetDB().Exec(query, username, email, githubID, avatar)
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
		Avatar:       avatar,
		AuthProvider: "github",
		IsVerified:   true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
