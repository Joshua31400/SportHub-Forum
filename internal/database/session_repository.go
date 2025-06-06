package database

import (
	"SportHub-Forum/internal/models"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

const (
	sessionCookieName = "session_token"
	sessionDuration   = 24 * time.Hour // 24h session duration
)

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func CreateSession(w http.ResponseWriter, userID int) error {
	// Call generateToken to create a new session token
	token, err := generateToken()
	if err != nil {
		return fmt.Errorf("Generate token", err)
	}

	// The token is valid for 24 hours call
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

func ValidateSession(r *http.Request) (int, bool) {
	// Get the session cookie from the request
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

	// Verify if the session is still valid
	if time.Now().After(expiresAt) {
		// If the session has expired, delete it
		DeleteSession(token)
		return 0, false
	}
	return userID, true
}

func DeleteSession(token string) error {
	query := "DELETE FROM session WHERE sessiontoken = ?"
	_, err := GetDB().Exec(query, token)
	return err
}

// Delelete the session and remove the cookie from the client (logout button)
func EndSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		DeleteSession(cookie.Value)
	}

	// Delete the session cookie for the client
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	return nil
}
