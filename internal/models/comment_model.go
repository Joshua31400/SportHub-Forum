package models

import "time"

type Comment struct {
	ID        int
	Content   string
	PostID    int
	UserID    int
	Username  string
	CreatedAt time.Time
}
