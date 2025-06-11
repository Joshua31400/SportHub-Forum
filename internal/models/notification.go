package models

import "time"

type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userid"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	IsRead    bool      `json:"isread"`
	PostID    Post      `json:"post"`
	CreatedAt time.Time `json:"createdat"`
}
