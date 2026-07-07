package signaling

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/voip-app/pkg/api"
)

type SignalingClient struct {
	conn     *websocket.Conn
	mu       sync.Mutex
	handlers map[string][]func(api.SignalingMessage)
}

func NewSignalingClient(serverURL string) (*SignalingClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		return nil, err
	}

	client := &SignalingClient{
		conn:     conn,
		handlers: make(map[string][]func(api.SignalingMessage)),
	}

	go client.readPump()

	return client, nil
}

func (c *SignalingClient) readPump() {
	defer c.conn.Close()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.trigger("close", api.SignalingMessage{Type: "close"})
			return
		}

		var msg api.SignalingMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		c.trigger(string(msg.Type), msg)
		c.trigger("*", msg)
	}
}

func (c *SignalingClient) Send(msg api.SignalingMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *SignalingClient) On(event string, handler func(api.SignalingMessage)) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers[event] = append(c.handlers[event], handler)
}

func (c *SignalingClient) Off(event string, handler func(api.SignalingMessage)) {
	c.mu.Lock()
	defer c.mu.Unlock()

	handlers := c.handlers[event]
	for i, h := range handlers {
		if &h == &handler {
			c.handlers[event] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

func (c *SignalingClient) trigger(event string, msg api.SignalingMessage) {
	c.mu.Lock()
	handlers := c.handlers[event]
	c.mu.Unlock()

	for _, h := range handlers {
		h(msg)
	}
}

func (c *SignalingClient) Close() error {
	return c.conn.Close()
}

func (c *SignalingClient) Join(room, userID string) error {
	return c.Send(api.SignalingMessage{
		Type:     api.MessageTypeJoin,
		Room:     room,
		SenderID: userID,
	})
}

func (c *SignalingClient) Leave(room, userID string) error {
	return c.Send(api.SignalingMessage{
		Type:     api.MessageTypeLeave,
		Room:     room,
		SenderID: userID,
	})
}

func (c *SignalingClient) SendOffer(targetID, sdp string) error {
	return c.Send(api.SignalingMessage{
		Type:     api.MessageTypeOffer,
		TargetID: targetID,
		SDP:      sdp,
	})
}

func (c *SignalingClient) SendAnswer(targetID, sdp string) error {
	return c.Send(api.SignalingMessage{
		Type:     api.MessageTypeAnswer,
		TargetID: targetID,
		SDP:      sdp,
	})
}

func (c *SignalingClient) SendICE(targetID, candidate string) error {
	return c.Send(api.SignalingMessage{
		Type:      api.MessageTypeICECandidate,
		TargetID:  targetID,
		Candidate: candidate,
	})
}

func (c *SignalingClient) SendTextMessage(channelID, content string) error {
	return c.Send(api.SignalingMessage{
		Type:      api.MessageTypeTextMessage,
		ChannelID: channelID,
		Content:   content,
	})
}

var _ = log.Printf
