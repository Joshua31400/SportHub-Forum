package models

import "time"

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userid"`
	Title     string    `json:"titre"`
	Content   string    `json:"contenu"`
	ImageURL  string    `json:"imageurl"`
	CreatedAt time.Time `json:"createdat"`
}

type PostCategory struct {
	PostID     int `json:"postid"`
	CategoryID int `json:"categoryid"`
}

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
