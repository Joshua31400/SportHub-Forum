package handlers

import (
	"SportHub-Forum/internal/authentification"
	"SportHub-Forum/internal/database"
	"SportHub-Forum/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// AuthHandler struct for Google OAuth authentication endpoints
type AuthHandler struct{}

// NewAuthHandler creates and returns a new AuthHandler instance
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// GoogleLogin initiates the Google OAuth2 flow by generating state and redirecting to Google
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := authentification.GenerateStateString()

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   false,
		MaxAge:   600,
	})

	url := authentification.GoogleConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handles Google OAuth callback, exchanges code for token and creates user session
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code missing", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := authentification.GoogleConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Error exchanging token", http.StatusInternalServerError)
		return
	}

	userInfo, err := h.getUserInfoFromGoogle(token.AccessToken)
	if err != nil {
		http.Error(w, "Error retrieving user information", http.StatusInternalServerError)
		return
	}

	user, err := h.processGoogleUser(userInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	err = database.CreateSession(w, user.UserID)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// getUserInfoFromGoogle fetches user profile information from Google API using access token
func (h *AuthHandler) getUserInfoFromGoogle(accessToken string) (*models.GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("error calling Google API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK response from Google: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Google response: %v", err)
	}

	var userInfo models.GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("error parsing Google JSON: %v", err)
	}

	return &userInfo, nil
}

// processGoogleUser handles Google user processing: finds existing user or creates new one
func (h *AuthHandler) processGoogleUser(googleUser *models.GoogleUserInfo) (*models.User, error) {
	// Check if user already exists by Google ID
	existingUser, err := database.GetUserByGoogleID(googleUser.ID)
	if err == nil && existingUser != nil {
		// Update last login for existing user
		database.UpdateUserLastLogin(existingUser.UserID)
		return existingUser, nil
	}

	// Check if user exists with same email (local account conflict)
	existingEmailUser, err := database.GetUserByEmail(googleUser.Email)
	if err == nil && existingEmailUser != nil {
		return nil, fmt.Errorf("local account with email %s already exists. Please login with your password", googleUser.Email)
	}

	// Create new Google user
	newUser, err := database.CreateGoogleUser(googleUser.Email, googleUser.Name, googleUser.ID, googleUser.Picture)
	if err != nil {
		return nil, fmt.Errorf("error creating Google user: %v", err)
	}

	return newUser, nil
}
