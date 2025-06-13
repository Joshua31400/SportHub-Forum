package models

import "time"

type User struct {
	UserID    int       `json:"userid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdat"`

	GitHubID     string    `json:"github_id,omitempty"`
	Avatar       string    `json:"avatar,omitempty"`
	IsVerified   bool      `json:"is_verified"`
	AuthProvider string    `json:"auth_provider"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type GitHubUserInfo struct {
	ID     int    `json:"id"`
	Login  string `json:"login"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"avatar_url"`
}
