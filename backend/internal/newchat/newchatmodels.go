// File: internal/newchat/newchatmodels.go

package newchat

import (
	"context"
	"scav/internal/media"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ------------------------- Types -------------------------

type Hub struct {
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan broadcastMsg

	mu       sync.Mutex
	stopped  bool
	stopChan chan struct{}
	stopOnce sync.Once
}

type Client struct {
	Conn   *websocket.Conn
	Send   chan []byte
	Room   string
	UserID string

	ctx    context.Context
	cancel context.CancelFunc
}

type Attachment struct {
	Filename string `db:"filename" json:"filename"`
	Path     string `db:"path" json:"path"`
}

type Message struct {
	ChatID     string       `db:"chatid"              json:"chatid"`
	UserID     string       `db:"sender"              json:"sender"`
	Text       string       `db:"text,omitempty" json:"text,omitempty"`
	FileURL    string       `db:"fileURL,omitempty" json:"fileURL,omitempty"`
	FileType   string       `db:"fileType,omitempty" json:"fileType,omitempty"` // "image" or "video"
	CreatedAt  time.Time    `db:"createdAt" json:"createdAt"`
	ReplyTo    *ReplyRef    `db:"replyTo,omitempty" json:"replyTo,omitempty"`
	SenderName string       `db:"senderName,omitempty" json:"senderName,omitempty"`
	AvatarURL  string       `db:"avatarUrl,omitempty"   json:"avatarUrl,omitempty"`
	Media      *media.Media `db:"media,omitempty"   json:"media,omitempty"`
	EditedAt   *time.Time   `db:"editedAt,omitempty" json:"editedAt,omitempty"`
	Deleted    bool         `db:"deleted"           json:"deleted"`
	ReadBy     []string     `db:"readBy,omitempty"  json:"readBy,omitempty"`
	Status     string       `db:"status,omitempty"  json:"status,omitempty"` // e.g. "sent", "read"

	MessageID string       `db:"messageid" json:"messageid"`
	Room      string       `db:"room" json:"room"`
	SenderID  string       `db:"senderid" json:"senderid"`
	Content   string       `db:"content" json:"content"`
	Files     []Attachment `db:"files,omitempty" json:"files,omitempty"`
	Timestamp int64        `db:"timestamp" json:"timestamp"`
}

type inboundPayload struct {
	Action  string `json:"action"`
	ID      string `json:"id,omitempty"`
	Content string `json:"content,omitempty"`
}

type outboundPayload struct {
	Action    string       `json:"action"`
	ID        string       `json:"id,omitempty"`
	Room      string       `json:"room,omitempty"`
	SenderID  string       `json:"senderid,omitempty"`
	Content   string       `json:"content,omitempty"`
	Files     []Attachment `json:"files,omitempty"`
	Timestamp int64        `json:"timestamp,omitempty"`
}

type broadcastMsg struct {
	Room string
	Data []byte
}

type Chat struct {
	Users        []string        `db:"users,omitempty"        json:"users,omitempty"`
	LastMessage  *MessagePreview `db:"lastMessage,omitempty"  json:"lastMessage,omitempty"`
	ReadStatus   map[string]bool `db:"readStatus,omitempty"   json:"readStatus,omitempty"`
	ChatID       string          `db:"chatid,omitempty"       json:"chatid"`
	Participants []string        `db:"participants,omitempty" json:"participants,omitempty"`
	CreatedAt    time.Time       `db:"createdAt"              json:"createdAt"`
	UpdatedAt    time.Time       `db:"updatedAt"              json:"updatedAt"`
	EntityType   string          `db:"entitytype,omitempty"   json:"entitytype,omitempty"`
	EntityId     string          `db:"entityid,omitempty"     json:"entityid,omitempty"`
	LastSeq      int64           `db:"lastSeq,omitempty"      json:"lastSeq,omitempty"`
}

type MessagePreview struct {
	Text      string    `db:"text"      json:"text"`
	UserID    string    `db:"userid"    json:"userid"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
}

// ReplyRef represents the client-side "replyTo" payload.
type ReplyRef struct {
	ID   string `json:"id"`
	User string `json:"user"`
	Text string `json:"text"`
}

// IncomingWSMessage represents a generic WebSocket inbound payload
type IncomingWSMessage struct {
	Type      string `json:"type"`
	ChatID    string `json:"chatid,omitempty"`
	MessageID string `json:"messageid,omitempty"`
	Content   string `json:"content,omitempty"`
	MediaURL  string `json:"mediaUrl,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
	Online    bool   `json:"online,omitempty"`
	ClientID  string `json:"clientId,omitempty"`
}

// type searchResult struct {
// 	Matches []ChatMessage `json:"matches"`
// }

// type chatMessage struct {
// 	ID        string `json:"id"`
// 	Sender    string `json:"sender"`
// 	Text      string `json:"text"`
// 	Timestamp string `json:"timestamp"`
// }
