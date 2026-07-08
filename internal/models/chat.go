package models

import "time"

type Chat struct {
	ID int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Title *string `json:"title"`
	LastMessage *string `json:"last_message"`
	CreatedAt time.Time `json:"created_at"`
}