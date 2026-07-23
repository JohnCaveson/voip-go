package channel

import (
	"errors"
	"strings"

	"github.com/voip-app/internal/models"
)

var (
	ErrNotFound      = errors.New("room not found")
	ErrDefaultDelete = errors.New("cannot delete default room")
	ErrDuplicateName = errors.New("room name already exists")
	ErrEmptyName     = errors.New("room name cannot be empty")
	ErrInvalidType   = errors.New("invalid room type")
)

type ChannelInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsDefault bool   `json:"is_default"`
}

func ParseType(s string) models.ChannelType {
	switch strings.ToLower(s) {
	case "text", "chat":
		return models.ChannelTypeText
	case "voice", "audio":
		return models.ChannelTypeVoice
	default:
		return models.ChannelType("")
	}
}
