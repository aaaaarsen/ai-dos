package models

import "time"

type Summary struct {
	ID int64 `json:"id"`
	ChatID int64 `json:"chat_id"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}