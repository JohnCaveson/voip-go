package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/voip-app/pkg/models"
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

func TestUpdateChannel(t *testing.T) {
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

	ch.Name = "#updated"
	ch.Type = models.ChannelTypeVoice
	if err := s.UpdateChannel(ctx, ch); err != nil {
		t.Fatal(err)
	}

	got, err := s.GetChannel(ctx, "ch1")
	if err != nil {
		t.Fatal(err)
	}

	if got.Name != "#updated" {
		t.Errorf("expected name #updated, got %s", got.Name)
	}
	if got.Type != models.ChannelTypeVoice {
		t.Errorf("expected type voice, got %s", got.Type)
	}
}

func TestListChannelsEmpty(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	list, err := s.ListChannels(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 0 {
		t.Errorf("expected 0 channels, got %d", len(list))
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

func TestListMessagesPagination(t *testing.T) {
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

	for i := 0; i < 10; i++ {
		msg := &models.Message{
			ID:        "m" + string(rune('0'+i)),
			ChannelID: "ch1",
			UserID:    "u1",
			Username:  "alice",
			Content:   "message",
			CreatedAt: time.Now().Add(time.Duration(i) * time.Second),
		}
		if err := s.SendMessage(ctx, msg); err != nil {
			t.Fatal(err)
		}
	}

	list, err := s.ListMessages(ctx, "ch1", 5, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 5 {
		t.Errorf("expected 5 messages with limit 5, got %d", len(list))
	}

	list2, err := s.ListMessages(ctx, "ch1", 5, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(list2) != 5 {
		t.Errorf("expected 5 messages with offset 5, got %d", len(list2))
	}

	if list[0].ID == list2[0].ID {
		t.Error("expected different messages for different offsets")
	}
}

func TestListMessagesDefaultLimit(t *testing.T) {
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

	msg := &models.Message{
		ID:        "m1",
		ChannelID: "ch1",
		UserID:    "u1",
		Username:  "alice",
		Content:   "hello",
		CreatedAt: time.Now(),
	}
	if err := s.SendMessage(ctx, msg); err != nil {
		t.Fatal(err)
	}

	list, err := s.ListMessages(ctx, "ch1", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 message, got %d", len(list))
	}
}

func TestListMessagesFromDifferentChannel(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	ch1 := &models.Channel{ID: "ch1", Name: "#general", Type: models.ChannelTypeText, CreatedAt: time.Now()}
	ch2 := &models.Channel{ID: "ch2", Name: "#random", Type: models.ChannelTypeText, CreatedAt: time.Now()}
	s.CreateChannel(ctx, ch1)
	s.CreateChannel(ctx, ch2)

	msg1 := &models.Message{ID: "m1", ChannelID: "ch1", UserID: "u1", Username: "alice", Content: "in ch1", CreatedAt: time.Now()}
	msg2 := &models.Message{ID: "m2", ChannelID: "ch2", UserID: "u1", Username: "alice", Content: "in ch2", CreatedAt: time.Now()}
	s.SendMessage(ctx, msg1)
	s.SendMessage(ctx, msg2)

	list, err := s.ListMessages(ctx, "ch1", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 message from ch1, got %d", len(list))
	}
	if list[0].Content != "in ch1" {
		t.Errorf("expected content 'in ch1', got %s", list[0].Content)
	}
}

func TestDeleteMessage(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	ch := &models.Channel{ID: "ch1", Name: "#general", Type: models.ChannelTypeText, CreatedAt: time.Now()}
	s.CreateChannel(ctx, ch)

	msg := &models.Message{ID: "m1", ChannelID: "ch1", UserID: "u1", Username: "alice", Content: "hello", CreatedAt: time.Now()}
	s.SendMessage(ctx, msg)

	if err := s.DeleteMessage(ctx, "m1"); err != nil {
		t.Fatal(err)
	}

	list, err := s.ListMessages(ctx, "ch1", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Errorf("expected 0 messages after delete, got %d", len(list))
	}
}

func TestDeleteMessageNotFound(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	err := s.DeleteMessage(ctx, "nonexistent")
	if err != nil {
		t.Errorf("expected no error for deleting nonexistent message, got %v", err)
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

func TestGetUserNotFound(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	_, err := s.GetUser(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}

func TestListUsersEmpty(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	list, err := s.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Errorf("expected 0 users, got %d", len(list))
	}
}

func TestMultipleUsers(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	users := []*models.User{
		{ID: "u1", Username: "alice", JoinedAt: time.Now()},
		{ID: "u2", Username: "bob", JoinedAt: time.Now().Add(time.Second)},
		{ID: "u3", Username: "charlie", JoinedAt: time.Now().Add(2 * time.Second)},
	}

	for _, user := range users {
		if err := s.AddUser(ctx, user); err != nil {
			t.Fatal(err)
		}
	}

	list, err := s.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 users, got %d", len(list))
	}
}

func TestSetUserOnlineToggle(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	user := &models.User{ID: "u1", Username: "alice", JoinedAt: time.Now(), IsOnline: false}
	s.AddUser(ctx, user)

	s.SetUserOnline(ctx, "u1", true)
	got, _ := s.GetUser(ctx, "u1")
	if !got.IsOnline {
		t.Error("expected is_online true after setting online")
	}

	s.SetUserOnline(ctx, "u1", false)
	got, _ = s.GetUser(ctx, "u1")
	if got.IsOnline {
		t.Error("expected is_online false after setting offline")
	}
}

func TestChannelCreatedAtPreserved(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	now := time.Now().Truncate(time.Second)
	ch := &models.Channel{
		ID:        "ch1",
		Name:      "#general",
		Type:      models.ChannelTypeText,
		CreatedAt: now,
	}
	s.CreateChannel(ctx, ch)

	got, _ := s.GetChannel(ctx, "ch1")
	if !got.CreatedAt.Equal(now) {
		t.Errorf("expected created_at %v, got %v", now, got.CreatedAt)
	}
}

func TestMessageCreatedAtPreserved(t *testing.T) {
	s := newTestDB(t)
	ctx := context.Background()

	ch := &models.Channel{ID: "ch1", Name: "#general", Type: models.ChannelTypeText, CreatedAt: time.Now()}
	s.CreateChannel(ctx, ch)

	now := time.Now().Truncate(time.Second)
	msg := &models.Message{ID: "m1", ChannelID: "ch1", UserID: "u1", Username: "alice", Content: "hello", CreatedAt: now}
	s.SendMessage(ctx, msg)

	list, _ := s.ListMessages(ctx, "ch1", 10, 0)
	if len(list) != 1 {
		t.Fatal("expected 1 message")
	}
	if !list[0].CreatedAt.Equal(now) {
		t.Errorf("expected created_at %v, got %v", now, list[0].CreatedAt)
	}
}
