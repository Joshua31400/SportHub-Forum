package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"postid"`
	UserID    int       `json:"userid"`
	Contenu   string    `json:"content"`
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
