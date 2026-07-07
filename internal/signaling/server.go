package signaling

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/voip-app/pkg/api"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	ID   string
	Room string
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]map[string]*Client),
	}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	go client.writePump()
	go h.readPump(client)
}

func (h *Hub) readPump(client *Client) {
	defer func() {
		room, userID := h.unregister(client)
		client.conn.Close()
		if room != "" {
			h.broadcastToRoom(room, api.SignalingMessage{
				Type:     api.MessageTypePeerLeft,
				SenderID: userID,
			}, "")
		}
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg api.SignalingMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case api.MessageTypeJoin:
			client.ID = msg.SenderID
			client.Room = msg.Room
			h.register(client)

			h.broadcastToRoom(client.Room, api.SignalingMessage{
				Type:     api.MessageTypePeerJoined,
				SenderID: client.ID,
			}, client.ID)

		case api.MessageTypeLeave:
			room, _ := h.unregister(client)
			if room != "" {
				h.broadcastToRoom(room, api.SignalingMessage{
					Type:     api.MessageTypePeerLeft,
					SenderID: client.ID,
				}, "")
			}

		case api.MessageTypeOffer, api.MessageTypeAnswer:
			h.sendToClient(msg.TargetID, message)

		case api.MessageTypeICECandidate:
			h.sendToClient(msg.TargetID, message)

		case api.MessageTypeTextMessage:
			h.broadcastToRoom(client.Room, msg, client.ID)
		}
	}
}

func (h *Hub) register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[client.Room] == nil {
		h.rooms[client.Room] = make(map[string]*Client)
	}
	h.rooms[client.Room][client.ID] = client
}

func (h *Hub) unregister(client *Client) (room string, userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room = client.Room
	userID = client.ID

	if clients, ok := h.rooms[room]; ok {
		delete(clients, userID)
		if len(clients) == 0 {
			delete(h.rooms, room)
		}
	}

	return
}

func (h *Hub) sendToClient(targetID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, clients := range h.rooms {
		if client, ok := clients[targetID]; ok {
			select {
			case client.send <- message:
			default:
			}
			return
		}
	}
}

func (h *Hub) broadcastToRoom(room string, msg api.SignalingMessage, excludeID string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	for id, client := range h.rooms[room] {
		if id != excludeID {
			select {
			case client.send <- data:
			default:
			}
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}
