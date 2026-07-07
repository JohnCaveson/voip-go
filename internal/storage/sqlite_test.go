package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/voip-app/internal/models"
)

func newTestDB(t *testing.T) *SQLiteStorage {
	t.Helper()

	f, err := os.CreateTemp("", "voip-test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	s, err := NewSQLiteStorage(f.Name())
	if err != nil {
		os.Remove(f.Name())
		t.Fatal(err)
	}

	t.Cleanup(func() {
		s.Close()
		os.Remove(f.Name())
	})

	return s
}

func TestCreateAndGetChannel(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	ch := &models.Channel{
		ID:        "ch1",
		Name:      "#general",
		Type:      models.ChannelTypeText,
		IsDefault: true,
		CreatedAt: time.Now(),
	}

	if err := s.CreateChannel(ctx, ch); err != nil {
		t.Fatal(err)
	}

	got, err := s.GetChannel(ctx, "ch1")
	if err != nil {
		t.Fatal(err)
	}

	if got.Name != "#general" {
		t.Errorf("expected name #general, got %s", got.Name)
	}
	if got.Type != models.ChannelTypeText {
		t.Errorf("expected type text, got %s", got.Type)
	}
	if !got.IsDefault {
		t.Error("expected is_default true")
	}
}

func TestGetChannelNotFound(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	_, err := s.GetChannel(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent channel")
	}
}

func TestListChannels(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	channels := []*models.Channel{
		{ID: "ch1", Name: "#general", Type: models.ChannelTypeText, CreatedAt: time.Now()},
		{ID: "ch2", Name: "#random", Type: models.ChannelTypeText, CreatedAt: time.Now().Add(time.Second)},
		{ID: "ch3", Name: "🔊 General", Type: models.ChannelTypeVoice, CreatedAt: time.Now().Add(2 * time.Second)},
	}

	for _, ch := range channels {
		if err := s.CreateChannel(ctx, ch); err != nil {
			t.Fatal(err)
		}
	}

	list, err := s.ListChannels(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 3 {
		t.Fatalf("expected 3 channels, got %d", len(list))
	}
}

func TestSoftDeleteChannel(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	ch := &models.Channel{
		ID:        "ch1",
		Name:      "#general",
		Type:      models.ChannelTypeText,
		CreatedAt: time.Now(),
	}
	if err := s.CreateChannel(ctx, ch); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteChannel(ctx, "ch1"); err != nil {
		t.Fatal(err)
	}

	_, err := s.GetChannel(ctx, "ch1")
	if err == nil {
		t.Error("expected error for deleted channel")
	}

	list, err := s.ListChannels(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Errorf("expected 0 channels after delete, got %d", len(list))
	}
}

func TestSendAndListMessages(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	ch := &models.Channel{
		ID:        "ch1",
		Name:      "#general",
		Type:      models.ChannelTypeText,
		CreatedAt: time.Now(),
	}
	if err := s.CreateChannel(ctx, ch); err != nil {
		t.Fatal(err)
	}

	msgs := []*models.Message{
		{ID: "m1", ChannelID: "ch1", UserID: "u1", Username: "alice", Content: "hello", CreatedAt: time.Now()},
		{ID: "m2", ChannelID: "ch1", UserID: "u2", Username: "bob", Content: "hi", CreatedAt: time.Now().Add(time.Second)},
	}
	for _, msg := range msgs {
		if err := s.SendMessage(ctx, msg); err != nil {
			t.Fatal(err)
		}
	}

	list, err := s.ListMessages(ctx, "ch1", 10, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(list))
	}
}

func TestUserCRUD(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	user := &models.User{
		ID:       "u1",
		Username: "alice",
		JoinedAt: time.Now(),
		IsOnline: true,
	}

	if err := s.AddUser(ctx, user); err != nil {
		t.Fatal(err)
	}

	got, err := s.GetUser(ctx, "u1")
	if err != nil {
		t.Fatal(err)
	}

	if got.Username != "alice" {
		t.Errorf("expected username alice, got %s", got.Username)
	}
	if !got.IsOnline {
		t.Error("expected is_online true")
	}

	if err := s.SetUserOnline(ctx, "u1", false); err != nil {
		t.Fatal(err)
	}

	got, _ = s.GetUser(ctx, "u1")
	if got.IsOnline {
		t.Error("expected is_online false after update")
	}

	list, err := s.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 user, got %d", len(list))
	}
}
