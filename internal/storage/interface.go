package storage

import (
	"context"

	"github.com/voip-app/internal/models"
)

type ChannelStore interface {
	CreateChannel(ctx context.Context, ch *models.Channel) error
	GetChannel(ctx context.Context, id string) (*models.Channel, error)
	ListChannels(ctx context.Context) ([]*models.Channel, error)
	UpdateChannel(ctx context.Context, ch *models.Channel) error
	DeleteChannel(ctx context.Context, id string) error
}

type MessageStore interface {
	SendMessage(ctx context.Context, msg *models.Message) error
	ListMessages(ctx context.Context, channelID string, limit, offset int) ([]*models.Message, error)
	DeleteMessage(ctx context.Context, id string) error
}

type UserStore interface {
	AddUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id string) (*models.User, error)
	ListUsers(ctx context.Context) ([]*models.User, error)
	SetUserOnline(ctx context.Context, id string, online bool) error
}

type Storage interface {
	ChannelStore
	MessageStore
	UserStore
	Close() error
}
