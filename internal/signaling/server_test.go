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

func TestHubJoinAndPeerNotification(t *testing.T) {
	_, server := newTestServer(t)

	conn1 := newTestWS(t, server.URL)
	conn2 := newTestWS(t, server.URL)

	sendWS(t, conn1, api.SignalingMessage{
		Type:     api.MessageTypeJoin,
		Room:     "test-room",
		SenderID: "user1",
	})

	time.Sleep(50 * time.Millisecond)

	sendWS(t, conn2, api.SignalingMessage{
		Type:     api.MessageTypeJoin,
		Room:     "test-room",
		SenderID: "user2",
	})

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

	sendWS(t, conn1, api.SignalingMessage{
		Type:     api.MessageTypeJoin,
		Room:     "test-room",
		SenderID: "user1",
	})
	sendWS(t, conn2, api.SignalingMessage{
		Type:     api.MessageTypeJoin,
		Room:     "test-room",
		SenderID: "user2",
	})

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

	sendWS(t, conn1, api.SignalingMessage{
		Type:     api.MessageTypeJoin,
		Room:     "test-room",
		SenderID: "user1",
	})
	sendWS(t, conn2, api.SignalingMessage{
		Type:     api.MessageTypeJoin,
		Room:     "test-room",
		SenderID: "user2",
	})

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

	sendWS(t, conn1, api.SignalingMessage{
		Type:     api.MessageTypeJoin,
		Room:     "test-room",
		SenderID: "user1",
	})
	sendWS(t, conn2, api.SignalingMessage{
		Type:     api.MessageTypeJoin,
		Room:     "test-room",
		SenderID: "user2",
	})

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

	sendWS(t, conn1, api.SignalingMessage{
		Type:     api.MessageTypeJoin,
		Room:     "test-room",
		SenderID: "user1",
	})
	sendWS(t, conn2, api.SignalingMessage{
		Type:     api.MessageTypeJoin,
		Room:     "test-room",
		SenderID: "user2",
	})

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
