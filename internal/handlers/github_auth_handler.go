package handlers

import (
	"SportHub-Forum/internal/authentification"
	"SportHub-Forum/internal/database"
	"SportHub-Forum/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// GitHubAuthHandler handles GitHub OAuth authentication
type GitHubAuthHandler struct{}

// NewGitHubAuthHandler creates a new GitHub auth handler
func NewGitHubAuthHandler() *GitHubAuthHandler {
	return &GitHubAuthHandler{}
}

// GitHubLogin initiates GitHub OAuth flow
func (h *GitHubAuthHandler) GitHubLogin(w http.ResponseWriter, r *http.Request) {
	state := authentification.GenerateStateString()

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state_github",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   600,
	})

	url := authentification.GitHubConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GitHubCallback handles GitHub OAuth callback
func (h *GitHubAuthHandler) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code missing", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := authentification.GitHubConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Error exchanging token", http.StatusInternalServerError)
		return
	}

	userInfo, err := h.getUserInfoFromGitHub(token.AccessToken)
	if err != nil {
		http.Error(w, "Error retrieving user information", http.StatusInternalServerError)
		return
	}

	user, err := h.processGitHubUser(userInfo)
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

// GitHubUserInfo represents GitHub user data
type GitHubUserInfo struct {
	ID     int    `json:"id"`
	Login  string `json:"login"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"avatar_url"`
}

// getUserInfoFromGitHub fetches user info from GitHub API
func (h *GitHubAuthHandler) getUserInfoFromGitHub(accessToken string) (*GitHubUserInfo, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling GitHub API: %v", err)
	}
	defer resp.Body.Close()

	var userInfo GitHubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("error parsing GitHub response: %v", err)
	}

	return &userInfo, nil
}

// processGitHubUser creates or retrieves GitHub user
func (h *GitHubAuthHandler) processGitHubUser(githubUser *GitHubUserInfo) (*models.User, error) {
	githubID := strconv.Itoa(githubUser.ID)

	existingUser, err := database.GetUserByGitHubID(githubID)
	if err == nil && existingUser != nil {
		return existingUser, nil
	}

	if githubUser.Email != "" {
		existingEmailUser, err := database.GetUserByEmail(githubUser.Email)
		if err == nil && existingEmailUser != nil {
			return nil, fmt.Errorf("account with email %s already exists", githubUser.Email)
		}
	}

	username := githubUser.Login
	if githubUser.Name != "" {
		username = githubUser.Name
	}

	newUser, err := database.CreateGitHubUser(githubUser.Email, username, githubID, githubUser.Avatar)
	if err != nil {
		return nil, fmt.Errorf("error creating GitHub user: %v", err)
	}

	return newUser, nil
}
