package channel

import (
	"context"
	"os"
	"testing"

	"github.com/voip-app/internal/models"
	"github.com/voip-app/internal/storage"
)

func newTestManager(t *testing.T) *Manager {
	t.Helper()

	f, err := os.CreateTemp("", "voip-test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	s, err := storage.NewSQLiteStorage(f.Name())
	if err != nil {
		os.Remove(f.Name())
		t.Fatal(err)
	}

	t.Cleanup(func() {
		s.Close()
		os.Remove(f.Name())
	})

	return NewManager(s)
}

func TestInitCreatesDefaults(t *testing.T) {
	m := newTestManager(t)
	ctx := context.Background()

	if err := m.Init(ctx); err != nil {
		t.Fatal(err)
	}

	channels, err := m.List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(channels) != 2 {
		t.Fatalf("expected 2 default channels, got %d", len(channels))
	}

	names := make(map[string]*models.Channel)
	for _, ch := range channels {
		names[ch.Name] = ch
	}

	if _, ok := names[DefaultTextChannelName]; !ok {
		t.Errorf("missing default text channel: %s", DefaultTextChannelName)
	}
	if _, ok := names[DefaultVoiceChannelName]; !ok {
		t.Errorf("missing default voice channel: %s", DefaultVoiceChannelName)
	}
}

func TestInitIdempotent(t *testing.T) {
	m := newTestManager(t)
	ctx := context.Background()

	if err := m.Init(ctx); err != nil {
		t.Fatal(err)
	}
	if err := m.Init(ctx); err != nil {
		t.Fatal(err)
	}

	channels, err := m.List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(channels) != 2 {
		t.Errorf("expected 2 channels after second init, got %d", len(channels))
	}
}

func TestCreateChannel(t *testing.T) {
	m := newTestManager(t)
	ctx := context.Background()

	if err := m.Init(ctx); err != nil {
		t.Fatal(err)
	}

	ch, err := m.Create(ctx, "#random", models.ChannelTypeText)
	if err != nil {
		t.Fatal(err)
	}

	if ch.Name != "#random" {
		t.Errorf("expected #random, got %s", ch.Name)
	}
	if ch.IsDefault {
		t.Error("new channel should not be default")
	}

	channels, _ := m.List(ctx)
	if len(channels) != 3 {
		t.Errorf("expected 3 channels, got %d", len(channels))
	}
}

func TestCreateDuplicateName(t *testing.T) {
	m := newTestManager(t)
	ctx := context.Background()

	if err := m.Init(ctx); err != nil {
		t.Fatal(err)
	}

	_, err := m.Create(ctx, DefaultTextChannelName, models.ChannelTypeText)
	if err != ErrDuplicateName {
		t.Errorf("expected ErrDuplicateName, got %v", err)
	}
}

func TestCreateEmptyName(t *testing.T) {
	m := newTestManager(t)
	ctx := context.Background()

	_, err := m.Create(ctx, "", models.ChannelTypeText)
	if err != ErrEmptyName {
		t.Errorf("expected ErrEmptyName, got %v", err)
	}
}

func TestCreateInvalidType(t *testing.T) {
	m := newTestManager(t)
	ctx := context.Background()

	_, err := m.Create(ctx, "#test", models.ChannelType("invalid"))
	if err != ErrInvalidType {
		t.Errorf("expected ErrInvalidType, got %v", err)
	}
}

func TestDeleteDefaultFails(t *testing.T) {
	m := newTestManager(t)
	ctx := context.Background()

	if err := m.Init(ctx); err != nil {
		t.Fatal(err)
	}

	channels, _ := m.List(ctx)
	for _, ch := range channels {
		if ch.IsDefault {
			err := m.Delete(ctx, ch.ID)
			if err != ErrDefaultDelete {
				t.Errorf("expected ErrDefaultDelete for default channel %s, got %v", ch.Name, err)
			}
		}
	}
}

func TestDeleteNonDefault(t *testing.T) {
	m := newTestManager(t)
	ctx := context.Background()

	if err := m.Init(ctx); err != nil {
		t.Fatal(err)
	}

	ch, err := m.Create(ctx, "#random", models.ChannelTypeText)
	if err != nil {
		t.Fatal(err)
	}

	if err := m.Delete(ctx, ch.ID); err != nil {
		t.Fatal(err)
	}

	_, err = m.Get(ctx, ch.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	m := newTestManager(t)
	ctx := context.Background()

	err := m.Delete(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestRename(t *testing.T) {
	m := newTestManager(t)
	ctx := context.Background()

	if err := m.Init(ctx); err != nil {
		t.Fatal(err)
	}

	channels, _ := m.List(ctx)
	textCh := channels[0]

	if err := m.Rename(ctx, textCh.ID, "#new-name"); err != nil {
		t.Fatal(err)
	}

	updated, _ := m.Get(ctx, textCh.ID)
	if updated.Name != "#new-name" {
		t.Errorf("expected #new-name, got %s", updated.Name)
	}
}

func TestRenameDuplicate(t *testing.T) {
	m := newTestManager(t)
	ctx := context.Background()

	if err := m.Init(ctx); err != nil {
		t.Fatal(err)
	}

	ch1, _ := m.Create(ctx, "#ch1", models.ChannelTypeText)
	m.Create(ctx, "#ch2", models.ChannelTypeText)

	err := m.Rename(ctx, ch1.ID, "#ch2")
	if err != ErrDuplicateName {
		t.Errorf("expected ErrDuplicateName, got %v", err)
	}
}

func TestRenameEmpty(t *testing.T) {
	m := newTestManager(t)
	ctx := context.Background()

	if err := m.Init(ctx); err != nil {
		t.Fatal(err)
	}

	channels, _ := m.List(ctx)
	err := m.Rename(ctx, channels[0].ID, "")
	if err != ErrEmptyName {
		t.Errorf("expected ErrEmptyName, got %v", err)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTestManager(t)
	ctx := context.Background()

	_, err := m.Get(ctx, "nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
