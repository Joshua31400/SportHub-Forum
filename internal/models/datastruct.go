package models

import "time"

type User struct {
	UserID    int       `json:"userid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdat"`
}

type Session struct {
	ID           int       `json:"id"`
	UserID       int       `json:"userid"`
	SessionToken string    `json:"sessiontoken"`
	ExpiresAt    time.Time `json:"expiresat"`
}

type Category struct {
	ID  int    `json:"id"`
	Nom string `json:"nom"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userid"`
	Titre     string    `json:"titre"`
	Contenu   string    `json:"contenu"`
	ImageURL  string    `json:"imageurl"`
	CreatedAt time.Time `json:"createdat"`
}

type PostCategory struct {
	PostID      int `json:"postid"`
	CategorieID int `json:"categorieid"`
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"postid"`
	UserID    int       `json:"userid"`
	Contenu   string    `json:"contenu"`
	CreatedAt time.Time `json:"createdat"`
}

type LikePost struct {
	UserID int    `json:"userid"`
	PostID int    `json:"postid"`
	Type   string `json:"type"`
}

type LikeComment struct {
	UserID    int    `json:"userid"`
	CommentID int    `json:"commentid"`
	Type      string `json:"type"`
}
