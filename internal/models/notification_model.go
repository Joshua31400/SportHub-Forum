package models

import (
	"time"
)

type Notification struct {
	ID         int       `json:"id"`
	UserID     int       `json:"userId"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
	SourceType string    `json:"sourceType,omitempty"`
	SourceID   int       `json:"sourceId,omitempty"`
}
