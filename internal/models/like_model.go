package models

import (
	"time"
)

type Like struct {
	ID        int       `json:"id" db:"id"`
	PostID    int       `json:"postid" db:"postid"`
	UserID    int       `json:"userid" db:"userid"`
	CreatedAt time.Time `json:"createdat" db:"createdat"`
}
