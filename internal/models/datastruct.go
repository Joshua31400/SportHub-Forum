package models

import "time"

type User struct {
	UserID    int       `json:"userID"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type Session struct {
	ID           int       `json:"id"`
	UserID       int       `json:"userID"`
	SessionToken string    `json:"sessionToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

type Category struct {
	ID  int    `json:"id"`
	Nom string `json:"nom"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userID"`
	Titre     string    `json:"titre"`
	Contenu   string    `json:"contenu"`
	ImageURL  string    `json:"imageURL"`
	CreatedAt time.Time `json:"createdAt"`
}

type PostCategory struct {
	PostID      int `json:"postID"`
	CategorieID int `json:"categorieID"`
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"postID"`
	UserID    int       `json:"userID"`
	Contenu   string    `json:"contenu"`
	CreatedAt time.Time `json:"createdAt"`
}

type LikePost struct {
	UserID int    `json:"userID"`
	PostID int    `json:"postID"`
	Type   string `json:"type"`
}

type LikeComment struct {
	UserID    int    `json:"userID"`
	CommentID int    `json:"commentID"`
	Type      string `json:"type"`
}
