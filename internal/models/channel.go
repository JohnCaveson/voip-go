package models

import "time"

type ChannelType string

const (
	ChannelTypeText  ChannelType = "text"
	ChannelTypeVoice ChannelType = "voice"
)

type Channel struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Type      ChannelType `json:"type"`
	IsDefault bool        `json:"is_default"`
	CreatedAt time.Time   `json:"created_at"`
	Deleted   bool        `json:"-"`
}
