package models

import "time"

type User struct {
	UserID    int       `json:"userid" db:"userID"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password,omitempty" db:"password"`
	CreatedAt time.Time `json:"createdat" db:"createdAt"`

	GoogleID     string    `json:"google_id,omitempty" db:"google_id"`
	GitHubID     string    `json:"github_id,omitempty" db:"github_id"`
	Avatar       string    `json:"avatar,omitempty" db:"avatar"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	AuthProvider string    `json:"auth_provider" db:"auth_provider"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type CreateUserRequest struct {
	Email        string `json:"email" validate:"required,email"`
	Username     string `json:"username" validate:"required,min=3"`
	Password     string `json:"password,omitempty"`
	GoogleID     string `json:"google_id,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	AuthProvider string `json:"auth_provider"`
}

type Session struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"userid" db:"userid"`
	SessionToken string    `json:"sessiontoken" db:"sessiontoken"`
	ExpiresAt    time.Time `json:"expiresat" db:"expiresat"`
}

type Category struct {
	ID  int    `json:"id" db:"id"`
	Nom string `json:"nom" db:"nom"`
}

type Post struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"userid" db:"userID"`
	Titre     string    `json:"titre" db:"titre"`
	Contenu   string    `json:"contenu" db:"contenu"`
	ImageURL  string    `json:"imageurl" db:"imageurl"`
	CreatedAt time.Time `json:"createdat" db:"createdat"`
}

type PostCategory struct {
	PostID      int `json:"postid" db:"postid"`
	CategorieID int `json:"categorieid" db:"categorieid"`
}

type Comment struct {
	ID        int       `json:"id" db:"id"`
	PostID    int       `json:"postid" db:"postid"`
	UserID    int       `json:"userid" db:"userID"`
	Contenu   string    `json:"contenu" db:"contenu"`
	CreatedAt time.Time `json:"createdat" db:"createdat"`
}

type LikePost struct {
	UserID int    `json:"userid" db:"userID"`
	PostID int    `json:"postid" db:"postid"`
	Type   string `json:"type" db:"type"`
}

type LikeComment struct {
	UserID    int    `json:"userid" db:"userID"`
	CommentID int    `json:"commentid" db:"commentid"`
	Type      string `json:"type" db:"type"`
}
