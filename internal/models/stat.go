package models

import (
	"time"
)

type DayStat struct {
    Day          time.Time `json:"day"`
    MessageCount int       `json:"message_count"`
}