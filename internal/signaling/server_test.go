package signaling

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/voip-app/pkg/api"
)

func newTestServer(t *testing.T) (*Hub, *httptest.Server) {
	t.Helper()

	hub := NewHub()
	server := httptest.NewServer(hub)
	t.Cleanup(server.Close)

	return hub, server
}

func newTestWS(t *testing.T, serverURL string) *websocket.Conn {
	t.Helper()

	wsURL := "ws" + strings.TrimPrefix(serverURL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })

	return conn
}

func sendWS(t *testing.T, conn *websocket.Conn, msg api.SignalingMessage) {
	t.Helper()

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		t.Fatal(err)
	}
}

func readWS(t *testing.T, conn *websocket.Conn) api.SignalingMessage {
	t.Helper()

	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}

	var msg api.SignalingMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatal(err)
	}
	return msg
}

func joinRoom(t *testing.T, conn *websocket.Conn, room, userID string) {
	t.Helper()
	sendWS(t, conn, api.SignalingMessage{
		Type:     api.MessageTypeJoin,
		Room:     room,
		SenderID: userID,
	})
	time.Sleep(50 * time.Millisecond)
}

func TestHubJoinAndPeerNotification(t *testing.T) {
	_, server := newTestServer(t)

	conn1 := newTestWS(t, server.URL)
	conn2 := newTestWS(t, server.URL)

	joinRoom(t, conn1, "test-room", "user1")
	joinRoom(t, conn2, "test-room", "user2")

	msg1 := readWS(t, conn1)
	if msg1.Type != api.MessageTypePeerJoined {
		t.Errorf("expected peer_joined, got %s", msg1.Type)
	}
	if msg1.SenderID != "user2" {
		t.Errorf("expected user2, got %s", msg1.SenderID)
	}
}

func TestHubOfferRelay(t *testing.T) {
	_, server := newTestServer(t)

	conn1 := newTestWS(t, server.URL)
	conn2 := newTestWS(t, server.URL)

	joinRoom(t, conn1, "test-room", "user1")
	joinRoom(t, conn2, "test-room", "user2")

	_ = readWS(t, conn1)

	sendWS(t, conn1, api.SignalingMessage{
		Type:     api.MessageTypeOffer,
		TargetID: "user2",
		SenderID: "user1",
		SDP:      "test-sdp",
	})

	msg := readWS(t, conn2)
	if msg.Type != api.MessageTypeOffer {
		t.Errorf("expected offer, got %s", msg.Type)
	}
	if msg.SDP != "test-sdp" {
		t.Errorf("expected test-sdp, got %s", msg.SDP)
	}
}

func TestHubICERelay(t *testing.T) {
	_, server := newTestServer(t)

	conn1 := newTestWS(t, server.URL)
	conn2 := newTestWS(t, server.URL)

	joinRoom(t, conn1, "test-room", "user1")
	joinRoom(t, conn2, "test-room", "user2")

	_ = readWS(t, conn1)

	sendWS(t, conn1, api.SignalingMessage{
		Type:      api.MessageTypeICECandidate,
		TargetID:  "user2",
		SenderID:  "user1",
		Candidate: "test-candidate",
	})

	msg := readWS(t, conn2)
	if msg.Type != api.MessageTypeICECandidate {
		t.Errorf("expected ice_candidate, got %s", msg.Type)
	}
	if msg.Candidate != "test-candidate" {
		t.Errorf("expected test-candidate, got %s", msg.Candidate)
	}
}

func TestHubTextMessageRelay(t *testing.T) {
	_, server := newTestServer(t)

	conn1 := newTestWS(t, server.URL)
	conn2 := newTestWS(t, server.URL)

	joinRoom(t, conn1, "test-room", "user1")
	joinRoom(t, conn2, "test-room", "user2")

	_ = readWS(t, conn1)

	sendWS(t, conn1, api.SignalingMessage{
		Type:      api.MessageTypeTextMessage,
		ChannelID: "ch1",
		SenderID:  "user1",
		Content:   "hello",
	})

	msg := readWS(t, conn2)
	if msg.Type != api.MessageTypeTextMessage {
		t.Errorf("expected text_message, got %s", msg.Type)
	}
	if msg.Content != "hello" {
		t.Errorf("expected hello, got %s", msg.Content)
	}
}

func TestHubLeave(t *testing.T) {
	_, server := newTestServer(t)

	conn1 := newTestWS(t, server.URL)
	conn2 := newTestWS(t, server.URL)

	joinRoom(t, conn1, "test-room", "user1")
	joinRoom(t, conn2, "test-room", "user2")

	_ = readWS(t, conn1)

	sendWS(t, conn1, api.SignalingMessage{
		Type:     api.MessageTypeLeave,
		Room:     "test-room",
		SenderID: "user1",
	})

	msg := readWS(t, conn2)
	if msg.Type != api.MessageTypePeerLeft {
		t.Errorf("expected peer_left, got %s", msg.Type)
	}
	if msg.SenderID != "user1" {
		t.Errorf("expected user1, got %s", msg.SenderID)
	}
}

func TestHubHTTP(t *testing.T) {
	hub, server := newTestServer(t)

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	hub.mu.RLock()
	roomCount := len(hub.rooms)
	hub.mu.RUnlock()

	if roomCount != 0 {
		t.Errorf("expected 0 rooms, got %d", roomCount)
	}
}

func TestHubMultipleRooms(t *testing.T) {
	_, server := newTestServer(t)

	conn1 := newTestWS(t, server.URL)
	conn2 := newTestWS(t, server.URL)
	conn3 := newTestWS(t, server.URL)

	joinRoom(t, conn1, "room1", "user1")
	joinRoom(t, conn2, "room2", "user2")
	joinRoom(t, conn3, "room1", "user3")

	msg := readWS(t, conn1)
	if msg.Type != api.MessageTypePeerJoined {
		t.Errorf("expected peer_joined, got %s", msg.Type)
	}
	if msg.SenderID != "user3" {
		t.Errorf("expected user3, got %s", msg.SenderID)
	}

	sendWS(t, conn1, api.SignalingMessage{
		Type:      api.MessageTypeTextMessage,
		ChannelID: "ch1",
		SenderID:  "user1",
		Content:   "hello room1",
	})

	msg3 := readWS(t, conn3)
	if msg3.Content != "hello room1" {
		t.Errorf("expected 'hello room1', got %s", msg3.Content)
	}

	conn2.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	_, _, err := conn2.ReadMessage()
	if err == nil {
		t.Error("expected timeout for user in different room")
	}
}

func TestHubDisconnectCleanup(t *testing.T) {
	hub, server := newTestServer(t)

	conn1 := newTestWS(t, server.URL)
	conn2 := newTestWS(t, server.URL)

	joinRoom(t, conn1, "test-room", "user1")
	joinRoom(t, conn2, "test-room", "user2")

	_ = readWS(t, conn1)

	conn1.Close()

	time.Sleep(100 * time.Millisecond)

	msg := readWS(t, conn2)
	if msg.Type != api.MessageTypePeerLeft {
		t.Errorf("expected peer_left after disconnect, got %s", msg.Type)
	}
	if msg.SenderID != "user1" {
		t.Errorf("expected user1, got %s", msg.SenderID)
	}

	hub.mu.RLock()
	roomClients := hub.rooms["test-room"]
	hub.mu.RUnlock()

	if len(roomClients) != 1 {
		t.Errorf("expected 1 client in room after disconnect, got %d", len(roomClients))
	}
}

func TestHubEmptyRoomCleanup(t *testing.T) {
	hub, server := newTestServer(t)

	conn1 := newTestWS(t, server.URL)

	joinRoom(t, conn1, "temp-room", "user1")

	conn1.Close()

	time.Sleep(100 * time.Millisecond)

	hub.mu.RLock()
	_, roomExists := hub.rooms["temp-room"]
	hub.mu.RUnlock()

	if roomExists {
		t.Error("expected room to be cleaned up after last user left")
	}
}

func TestHubAnswerRelay(t *testing.T) {
	_, server := newTestServer(t)

	conn1 := newTestWS(t, server.URL)
	conn2 := newTestWS(t, server.URL)

	joinRoom(t, conn1, "test-room", "user1")
	joinRoom(t, conn2, "test-room", "user2")

	_ = readWS(t, conn1)

	sendWS(t, conn2, api.SignalingMessage{
		Type:     api.MessageTypeAnswer,
		TargetID: "user1",
		SenderID: "user2",
		SDP:      "answer-sdp",
	})

	msg := readWS(t, conn1)
	if msg.Type != api.MessageTypeAnswer {
		t.Errorf("expected answer, got %s", msg.Type)
	}
	if msg.SDP != "answer-sdp" {
		t.Errorf("expected answer-sdp, got %s", msg.SDP)
	}
}

func TestHubInvalidJSON(t *testing.T) {
	_, server := newTestServer(t)

	conn := newTestWS(t, server.URL)

	if err := conn.WriteMessage(websocket.TextMessage, []byte("invalid json")); err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"ping"}`)); err != nil {
		t.Fatal(err)
	}
}
