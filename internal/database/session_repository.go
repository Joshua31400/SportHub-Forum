package database

import (
	"SportHub-Forum/internal/authentification"
	"SportHub-Forum/internal/models"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	sessionCookieName = "session_token"
	sessionDuration   = 24 * time.Hour // 24h session duration
)

// Generates a secure random session token encoded in base64.
func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Creates a new session for a user, stores it in the database, and sets the session cookie.
func CreateSession(w http.ResponseWriter, userID int) error {
	token, err := generateToken()
	if err != nil {
		return fmt.Errorf("Generate token", err)
	}

	expiresAt := time.Now().Add(sessionDuration)

	session := models.Session{
		UserID:       userID,
		SessionToken: token,
		ExpiresAt:    expiresAt,
	}

	db := GetDB()
	if db == nil {
		return fmt.Errorf("Connection with the database is not initialized")
	}

	query := "INSERT INTO session (userid, sessiontoken, expiresat) VALUES (?, ?, ?)"
	_, err = db.Exec(query, session.UserID, session.SessionToken, expiresAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("Insert in the databese: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
	})
	return nil
}

// Validates a session from the request cookie and checks if it is still valid.
func ValidateSession(r *http.Request) (int, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return 0, false
	}

	token := cookie.Value

	var userID int
	var expiresAtStr string

	query := "SELECT userid, expiresat FROM session WHERE sessiontoken = ?"
	err = GetDB().QueryRow(query, token).Scan(&userID, &expiresAtStr)
	if err != nil {
		return 0, false
	}

	expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil {
		return 0, false
	}

	if time.Now().After(expiresAt) {
		DeleteSession(token)
		return 0, false
	}
	return userID, true
}

// Ends the current session by deleting it and removing the session cookie.
func EndSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		token := cookie.Value
		DeleteSession(token)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// Deletes a session from the database using the session token.
func DeleteSession(token string) error {
	if token == "" {
		return fmt.Errorf("token vide")
	}

	db := GetDB()
	if db == nil {
		return fmt.Errorf("connexion DB non initialisée")
	}

	query := "DELETE FROM session WHERE sessiontoken = ?"
	_, err := db.Exec(query, token)
	return err
}

// Retrieves a user from the database by their Google ID.
func GetUserByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	var createdAtStr string
	var password sql.NullString
	var avatar sql.NullString
	var authProvider sql.NullString
	var isVerified sql.NullBool
	var updatedAtStr sql.NullString

	query := `SELECT userID, username, email, password, 
                 DATE_FORMAT(createdAt, '%Y-%m-%d %H:%i:%s') as createdAt,
                 google_id, avatar, auth_provider, is_verified,
                 DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') as updated_at
          FROM user WHERE google_id = ?`
	err := GetDB().QueryRow(query, googleID).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&password,
		&createdAtStr,
		&user.GoogleID,
		&avatar,
		&authProvider,
		&isVerified,
		&updatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving user by Google ID: %v", err)
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

// Creates a new user with Google authentication and returns the created user.
func CreateGoogleUser(email, username, googleID, avatar string) (*models.User, error) {
	now := time.Now()

	query := `
        INSERT INTO user (username, email, google_id, avatar, auth_provider, is_verified, createdat, updated_at)
        VALUES (?, ?, ?, ?, 'google', true, ?, ?)
    `

	result, err := ExecWithTimeout(query, username, email, googleID, avatar, now, now)
	if err != nil {
		return nil, fmt.Errorf("error creating Google user: %v", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return nil, fmt.Errorf("no rows affected, user not created")
	}

	return GetUserByGoogleID(googleID)
}

// Updates the last login timestamp for a user.
func UpdateUserLastLogin(userID int) error {
	query := `UPDATE user SET updated_at = ? WHERE userid = ?`
	_, err := ExecWithTimeout(query, time.Now(), userID)
	return err
}

// Authenticates a user by username and password.
func AuthenticateUser(username, password string) (*models.User, error) {
	user, err := GetUserByEmail(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("utilisateur non trouvé")
	}

	if !authentification.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("mot de passe incorrect")
	}

	return user, nil
}
