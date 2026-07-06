package models

import (
	"time"
)

type Mood struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    Emoji     string    `json:"emoji"`
    Date      time.Time `json:"date"`
    CreatedAt time.Time `json:"created_at"`
}

