package models

import "time"

type Session struct {
	ID           int       `json:"id"`
	UserID       int       `json:"userid"`
	SessionToken string    `json:"sessiontoken"`
	ExpiresAt    time.Time `json:"expiresat"`
}
