package channel

import (
	"errors"
	"strings"

	"github.com/voip-app/internal/models"
)

var (
	ErrNotFound      = errors.New("channel not found")
	ErrDefaultDelete = errors.New("cannot delete default channel")
	ErrDuplicateName = errors.New("channel name already exists")
	ErrEmptyName     = errors.New("channel name cannot be empty")
	ErrInvalidType   = errors.New("invalid channel type")
)

type ChannelInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsDefault bool   `json:"is_default"`
}

func ParseType(s string) models.ChannelType {
	switch strings.ToLower(s) {
	case "text":
		return models.ChannelTypeText
	case "voice":
		return models.ChannelTypeVoice
	default:
		return models.ChannelType("")
	}
}
