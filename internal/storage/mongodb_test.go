package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/voip-app/pkg/models"
)

func skipIfNoMongo(t *testing.T) {
	t.Helper()
	if os.Getenv("MONGODB_URI") == "" {
		t.Skip("MONGODB_URI not set, skipping MongoDB integration tests")
	}
}

func newTestMongoStorage(t *testing.T) *MongoDBStorage {
	t.Helper()
	uri := os.Getenv("MONGODB_URI")
	s, err := NewMongoDBStorage(uri)
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}
	t.Cleanup(func() {
		s.Close()
	})
	return s
}

func cleanupCollection(t *testing.T, s *MongoDBStorage, name string) {
	t.Helper()
	s.database.Collection(name).Drop(context.Background())
}

func TestMongoDBStorageImplementsInterface(t *testing.T) {
	var _ Storage = (*MongoDBStorage)(nil)
}

func TestMongoDBConnectionFailure(t *testing.T) {
	_, err := NewMongoDBStorage("mongodb://invalid-host:27017")
	if err == nil {
		t.Error("expected error for invalid MongoDB URI")
	}
}

func TestMongoDBChannelCRUD(t *testing.T) {
	skipIfNoMongo(t)
	s := newTestMongoStorage(t)
	cleanupCollection(t, s, "channels")
	ctx := context.Background()

	ch := &models.Channel{
		ID:        "test-ch-1",
		Name:      "Test Channel",
		Type:      models.ChannelTypeText,
		IsDefault: false,
		CreatedAt: time.Now(),
	}

	if err := s.CreateChannel(ctx, ch); err != nil {
		t.Fatalf("CreateChannel: %v", err)
	}

	got, err := s.GetChannel(ctx, "test-ch-1")
	if err != nil {
		t.Fatalf("GetChannel: %v", err)
	}
	if got.Name != "Test Channel" {
		t.Errorf("expected name 'Test Channel', got %s", got.Name)
	}
	if got.Type != models.ChannelTypeText {
		t.Errorf("expected type text, got %s", got.Type)
	}

	channels, err := s.ListChannels(ctx)
	if err != nil {
		t.Fatalf("ListChannels: %v", err)
	}
	if len(channels) != 1 {
		t.Errorf("expected 1 channel, got %d", len(channels))
	}

	got.Name = "Updated Channel"
	if err := s.UpdateChannel(ctx, got); err != nil {
		t.Fatalf("UpdateChannel: %v", err)
	}
	updated, _ := s.GetChannel(ctx, "test-ch-1")
	if updated.Name != "Updated Channel" {
		t.Errorf("expected updated name, got %s", updated.Name)
	}

	if err := s.DeleteChannel(ctx, "test-ch-1"); err != nil {
		t.Fatalf("DeleteChannel: %v", err)
	}
	deleted, err := s.GetChannel(ctx, "test-ch-1")
	if err == nil {
		t.Error("expected error for deleted channel")
	}
	_ = deleted
}

func TestMongoDBMessageCRUD(t *testing.T) {
	skipIfNoMongo(t)
	s := newTestMongoStorage(t)
	cleanupCollection(t, s, "messages")
	ctx := context.Background()

	msg := &models.Message{
		ID:        "test-msg-1",
		ChannelID: "ch-1",
		UserID:    "user-1",
		Username:  "tester",
		Content:   "hello world",
		CreatedAt: time.Now(),
	}

	if err := s.SendMessage(ctx, msg); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}

	messages, err := s.ListMessages(ctx, "ch-1", 50, 0)
	if err != nil {
		t.Fatalf("ListMessages: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}
	if messages[0].Content != "hello world" {
		t.Errorf("expected content 'hello world', got %s", messages[0].Content)
	}
	if messages[0].Username != "tester" {
		t.Errorf("expected username 'tester', got %s", messages[0].Username)
	}

	msg2 := &models.Message{
		ID:        "test-msg-2",
		ChannelID: "ch-1",
		UserID:    "user-2",
		Username:  "other",
		Content:   "second message",
		CreatedAt: time.Now().Add(time.Minute),
	}
	s.SendMessage(ctx, msg2)

	messages, _ = s.ListMessages(ctx, "ch-1", 1, 0)
	if len(messages) != 1 {
		t.Errorf("expected 1 message with limit 1, got %d", len(messages))
	}

	if err := s.DeleteMessage(ctx, "test-msg-1"); err != nil {
		t.Fatalf("DeleteMessage: %v", err)
	}
	messages, _ = s.ListMessages(ctx, "ch-1", 50, 0)
	if len(messages) != 1 {
		t.Errorf("expected 1 message after delete, got %d", len(messages))
	}
}

func TestMongoDBUserCRUD(t *testing.T) {
	skipIfNoMongo(t)
	s := newTestMongoStorage(t)
	cleanupCollection(t, s, "users")
	ctx := context.Background()

	user := &models.User{
		ID:       "test-user-1",
		Username: "alice",
		JoinedAt: time.Now(),
		IsOnline: true,
	}

	if err := s.AddUser(ctx, user); err != nil {
		t.Fatalf("AddUser: %v", err)
	}

	got, err := s.GetUser(ctx, "test-user-1")
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got.Username != "alice" {
		t.Errorf("expected username 'alice', got %s", got.Username)
	}
	if !got.IsOnline {
		t.Error("expected user to be online")
	}

	users, err := s.ListUsers(ctx)
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}

	if err := s.SetUserOnline(ctx, "test-user-1", false); err != nil {
		t.Fatalf("SetUserOnline: %v", err)
	}
	offline, _ := s.GetUser(ctx, "test-user-1")
	if offline.IsOnline {
		t.Error("expected user to be offline")
	}
}

func TestMongoDBDuplicateChannelName(t *testing.T) {
	skipIfNoMongo(t)
	s := newTestMongoStorage(t)
	cleanupCollection(t, s, "channels")
	ctx := context.Background()

	ch1 := &models.Channel{ID: "dup-1", Name: "Unique Name", Type: models.ChannelTypeText, CreatedAt: time.Now()}
	ch2 := &models.Channel{ID: "dup-2", Name: "Unique Name", Type: models.ChannelTypeVoice, CreatedAt: time.Now()}

	if err := s.CreateChannel(ctx, ch1); err != nil {
		t.Fatalf("CreateChannel 1: %v", err)
	}
	if err := s.CreateChannel(ctx, ch2); err == nil {
		t.Error("expected error for duplicate channel name")
	}
}

func TestMongoDBSoftDeleteFiltering(t *testing.T) {
	skipIfNoMongo(t)
	s := newTestMongoStorage(t)
	cleanupCollection(t, s, "channels")
	ctx := context.Background()

	ch := &models.Channel{ID: "soft-1", Name: "Soft Delete Test", Type: models.ChannelTypeText, CreatedAt: time.Now()}
	s.CreateChannel(ctx, ch)
	s.DeleteChannel(ctx, "soft-1")

	channels, _ := s.ListChannels(ctx)
	if len(channels) != 0 {
		t.Errorf("expected 0 channels after soft delete, got %d", len(channels))
	}

	_, err := s.GetChannel(ctx, "soft-1")
	if err == nil {
		t.Error("expected error for soft-deleted channel")
	}
}
