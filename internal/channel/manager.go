package channel

import (
	"context"
	"fmt"
	"time"

	"github.com/voip-app/internal/models"
	"github.com/voip-app/internal/storage"
)

const (
	DefaultChatRoomName  = "General"
	DefaultAudioRoomName = "Lounge"
)

type Manager struct {
	store storage.Storage
}

func NewManager(store storage.Storage) *Manager {
	return &Manager{store: store}
}

func (m *Manager) Init(ctx context.Context) error {
	channels, err := m.store.ListChannels(ctx)
	if err != nil {
		return fmt.Errorf("list rooms: %w", err)
	}

	existing := make(map[string]bool)
	for _, ch := range channels {
		existing[ch.Name] = true
	}

	defaults := []*models.Channel{
		{
			ID:        "default-text",
			Name:      DefaultChatRoomName,
			Type:      models.ChannelTypeText,
			IsDefault: true,
			CreatedAt: time.Now(),
		},
		{
			ID:        "default-voice",
			Name:      DefaultAudioRoomName,
			Type:      models.ChannelTypeVoice,
			IsDefault: true,
			CreatedAt: time.Now(),
		},
	}

	for _, ch := range defaults {
		if !existing[ch.Name] {
			if err := m.store.CreateChannel(ctx, ch); err != nil {
				return fmt.Errorf("create default room %s: %w", ch.Name, err)
			}
		}
	}

	return nil
}

func (m *Manager) Create(ctx context.Context, name string, chType models.ChannelType) (*models.Channel, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	if chType != models.ChannelTypeText && chType != models.ChannelTypeVoice {
		return nil, ErrInvalidType
	}

	channels, err := m.store.ListChannels(ctx)
	if err != nil {
		return nil, fmt.Errorf("list rooms: %w", err)
	}

	for _, ch := range channels {
		if ch.Name == name {
			return nil, ErrDuplicateName
		}
	}

	ch := &models.Channel{
		ID:        fmt.Sprintf("ch-%d", time.Now().UnixNano()),
		Name:      name,
		Type:      chType,
		IsDefault: false,
		CreatedAt: time.Now(),
	}

	if err := m.store.CreateChannel(ctx, ch); err != nil {
		return nil, fmt.Errorf("create room: %w", err)
	}

	return ch, nil
}

func (m *Manager) Get(ctx context.Context, id string) (*models.Channel, error) {
	ch, err := m.store.GetChannel(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	return ch, nil
}

func (m *Manager) List(ctx context.Context) ([]*models.Channel, error) {
	return m.store.ListChannels(ctx)
}

func (m *Manager) Delete(ctx context.Context, id string) error {
	ch, err := m.store.GetChannel(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	if ch.IsDefault {
		return ErrDefaultDelete
	}

	return m.store.DeleteChannel(ctx, id)
}

func (m *Manager) Rename(ctx context.Context, id, newName string) error {
	if newName == "" {
		return ErrEmptyName
	}

	ch, err := m.store.GetChannel(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	channels, err := m.store.ListChannels(ctx)
	if err != nil {
		return fmt.Errorf("list rooms: %w", err)
	}

	for _, other := range channels {
		if other.ID != id && other.Name == newName {
			return ErrDuplicateName
		}
	}

	ch.Name = newName
	return m.store.UpdateChannel(ctx, ch)
}
