package api

import (
	"testing"
	"time"

	"github.com/voip-app/pkg/models"
)

func TestUserToInfo(t *testing.T) {
	user := &models.User{
		ID:       "u1",
		Username: "alice",
		JoinedAt: time.Now(),
		IsOnline: true,
	}

	info := UserToInfo(user)

	if info.ID != "u1" {
		t.Errorf("expected ID u1, got %s", info.ID)
	}
	if info.Username != "alice" {
		t.Errorf("expected username alice, got %s", info.Username)
	}
}

func TestUserToInfoEmptyUser(t *testing.T) {
	user := &models.User{}
	info := UserToInfo(user)

	if info.ID != "" {
		t.Errorf("expected empty ID, got %s", info.ID)
	}
	if info.Username != "" {
		t.Errorf("expected empty username, got %s", info.Username)
	}
}

func TestMessageTypes(t *testing.T) {
	tests := []struct {
		msgType MessageType
		expected string
	}{
		{MessageTypeJoin, "join"},
		{MessageTypeLeave, "leave"},
		{MessageTypeOffer, "offer"},
		{MessageTypeAnswer, "answer"},
		{MessageTypeICECandidate, "ice_candidate"},
		{MessageTypeTextMessage, "text_message"},
		{MessageTypeScreenShare, "screen_share"},
		{MessageTypePeerJoined, "peer_joined"},
		{MessageTypePeerLeft, "peer_left"},
		{MessageTypeError, "error"},
	}

	for _, tt := range tests {
		if string(tt.msgType) != tt.expected {
			t.Errorf("MessageType %v: expected %s, got %s", tt.msgType, tt.expected, string(tt.msgType))
		}
	}
}

func TestSignalingMessageJSON(t *testing.T) {
	msg := SignalingMessage{
		Type:     MessageTypeJoin,
		Room:     "test-room",
		SenderID: "user1",
		TargetID: "user2",
		User: &UserInfo{
			ID:       "u1",
			Username: "alice",
		},
		SDP:      "test-sdp",
		Candidate: "test-candidate",
		ChannelID: "ch1",
		Content:   "hello",
		Error:     "test-error",
	}

	if msg.Type != MessageTypeJoin {
		t.Errorf("expected type join, got %s", msg.Type)
	}
	if msg.Room != "test-room" {
		t.Errorf("expected room test-room, got %s", msg.Room)
	}
	if msg.SenderID != "user1" {
		t.Errorf("expected sender_id user1, got %s", msg.SenderID)
	}
	if msg.TargetID != "user2" {
		t.Errorf("expected target_id user2, got %s", msg.TargetID)
	}
	if msg.User == nil || msg.User.ID != "u1" {
		t.Error("expected user with ID u1")
	}
	if msg.SDP != "test-sdp" {
		t.Errorf("expected SDP test-sdp, got %s", msg.SDP)
	}
	if msg.Candidate != "test-candidate" {
		t.Errorf("expected candidate test-candidate, got %s", msg.Candidate)
	}
	if msg.ChannelID != "ch1" {
		t.Errorf("expected channel_id ch1, got %s", msg.ChannelID)
	}
	if msg.Content != "hello" {
		t.Errorf("expected content hello, got %s", msg.Content)
	}
	if msg.Error != "test-error" {
		t.Errorf("expected error test-error, got %s", msg.Error)
	}
}

func TestRoomInfo(t *testing.T) {
	room := RoomInfo{
		Name: "test-room",
		Users: []UserInfo{
			{ID: "u1", Username: "alice"},
			{ID: "u2", Username: "bob"},
		},
	}

	if room.Name != "test-room" {
		t.Errorf("expected name test-room, got %s", room.Name)
	}
	if len(room.Users) != 2 {
		t.Errorf("expected 2 users, got %d", len(room.Users))
	}
	if room.Users[0].Username != "alice" {
		t.Errorf("expected first user alice, got %s", room.Users[0].Username)
	}
	if room.Users[1].Username != "bob" {
		t.Errorf("expected second user bob, got %s", room.Users[1].Username)
	}
}

func TestUserInfo(t *testing.T) {
	info := UserInfo{
		ID:       "u1",
		Username: "alice",
	}

	if info.ID != "u1" {
		t.Errorf("expected ID u1, got %s", info.ID)
	}
	if info.Username != "alice" {
		t.Errorf("expected username alice, got %s", info.Username)
	}
}
