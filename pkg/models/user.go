package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	JoinedAt  time.Time `json:"joined_at"`
	IsOnline  bool      `json:"is_online"`
}
