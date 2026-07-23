package models

import "time"

type Message struct {
	ID        string    `json:"id"`
	ChannelID string    `json:"channel_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
