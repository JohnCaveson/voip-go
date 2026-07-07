package signaling

import (
	"strings"
	"testing"
	"time"

	"github.com/voip-app/pkg/api"
)

func TestSignalingClientJoinAndReceive(t *testing.T) {
	_, server := newTestServer(t)

	client1, err := NewSignalingClient("ws" + strings.TrimPrefix(server.URL, "http"))
	if err != nil {
		t.Fatal(err)
	}
	defer client1.Close()

	client2, err := NewSignalingClient("ws" + strings.TrimPrefix(server.URL, "http"))
	if err != nil {
		t.Fatal(err)
	}
	defer client2.Close()

	peerJoined := make(chan string, 1)
	client1.On(string(api.MessageTypePeerJoined), func(msg api.SignalingMessage) {
		peerJoined <- msg.SenderID
	})

	if err := client1.Join("test-room", "user1"); err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	if err := client2.Join("test-room", "user2"); err != nil {
		t.Fatal(err)
	}

	select {
	case userID := <-peerJoined:
		if userID != "user2" {
			t.Errorf("expected user2, got %s", userID)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for peer_joined")
	}
}

func TestSignalingClientOffer(t *testing.T) {
	_, server := newTestServer(t)

	client1, err := NewSignalingClient("ws" + strings.TrimPrefix(server.URL, "http"))
	if err != nil {
		t.Fatal(err)
	}
	defer client1.Close()

	client2, err := NewSignalingClient("ws" + strings.TrimPrefix(server.URL, "http"))
	if err != nil {
		t.Fatal(err)
	}
	defer client2.Close()

	offerCh := make(chan api.SignalingMessage, 1)
	client2.On(string(api.MessageTypeOffer), func(msg api.SignalingMessage) {
		offerCh <- msg
	})

	client1.Join("test-room", "user1")
	time.Sleep(50 * time.Millisecond)
	client2.Join("test-room", "user2")

	time.Sleep(50 * time.Millisecond)

	if err := client1.SendOffer("user2", "test-sdp"); err != nil {
		t.Fatal(err)
	}

	select {
	case msg := <-offerCh:
		if msg.SDP != "test-sdp" {
			t.Errorf("expected test-sdp, got %s", msg.SDP)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for offer")
	}
}

func TestSignalingClientClose(t *testing.T) {
	_, server := newTestServer(t)

	client, err := NewSignalingClient("ws" + strings.TrimPrefix(server.URL, "http"))
	if err != nil {
		t.Fatal(err)
	}

	closeCh := make(chan bool, 1)
	client.On("close", func(msg api.SignalingMessage) {
		closeCh <- true
	})

	client.Close()

	select {
	case <-closeCh:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for close")
	}
}
