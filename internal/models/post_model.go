package models

import "time"

type Post struct {
	ID           int       `json:"id"`
	UserID       int       `json:"userid"`
	Title        string    `json:"titre"`
	Content      string    `json:"contenu"`
	ImageURL     string    `json:"imageurl"`
	CreatedAt    time.Time `json:"createdat"`
	CategoryID   int       `json:"categoryid"`
	CategoryName string    `json:"categoryname"`
	Username     string    `json:"username"`
}

// Future many-to-many relationship between posts and categories for multiple categories per post.gohtml
type PostCategory struct {
	PostID     int `json:"postid"`
	CategoryID int `json:"categoryid"`
}
