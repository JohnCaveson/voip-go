package api

import "github.com/voip-app/internal/models"

type MessageType string

const (
	MessageTypeJoin         MessageType = "join"
	MessageTypeLeave        MessageType = "leave"
	MessageTypeOffer        MessageType = "offer"
	MessageTypeAnswer       MessageType = "answer"
	MessageTypeICECandidate MessageType = "ice_candidate"
	MessageTypeTextMessage  MessageType = "text_message"
	MessageTypeScreenShare  MessageType = "screen_share"
	MessageTypePeerJoined   MessageType = "peer_joined"
	MessageTypePeerLeft     MessageType = "peer_left"
	MessageTypeError        MessageType = "error"
)

type SignalingMessage struct {
	Type     MessageType `json:"type"`
	Room     string      `json:"room,omitempty"`
	SenderID string      `json:"sender_id,omitempty"`
	TargetID string      `json:"target_id,omitempty"`
	User     *UserInfo   `json:"user,omitempty"`
	SDP      string      `json:"sdp,omitempty"`
	Candidate string     `json:"candidate,omitempty"`
	ChannelID string     `json:"channel_id,omitempty"`
	Content  string      `json:"content,omitempty"`
	Error    string      `json:"error,omitempty"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type RoomInfo struct {
	Name  string     `json:"name"`
	Users []UserInfo `json:"users"`
}

func UserToInfo(user *models.User) UserInfo {
	return UserInfo{
		ID:       user.ID,
		Username: user.Username,
	}
}
