// File: internal/mechat/mechatModels.go

package mechat

import (
	"scav/internal/verticals/media"
	"time"
)

type Message struct {
	MessageID  string       `db:"messageid,omitempty"  json:"messageid"`
	ChatID     string       `db:"chatid"               json:"chatid"`
	RoomID     string       `db:"roomid,omitempty"     json:"roomid,omitempty"`
	UserID     string       `db:"userid"               json:"userid"`
	Text       string       `db:"text,omitempty"       json:"text,omitempty"`
	FileURL    string       `db:"fileURL,omitempty"    json:"fileURL,omitempty"`
	FileType   string       `db:"fileType,omitempty"   json:"fileType,omitempty"` // "image" or "video"
	CreatedAt  time.Time    `db:"createdAt"            json:"createdAt"`
	ReplyTo    *ReplyRef    `db:"replyTo,omitempty"    json:"replyTo,omitempty"`
	SenderName string       `db:"senderName,omitempty" json:"senderName,omitempty"`
	AvatarURL  string       `db:"avatarUrl,omitempty"  json:"avatarUrl,omitempty"`
	Content    string       `db:"content"              json:"content"`
	Media      *media.Media `db:"media,omitempty"      json:"media,omitempty"`
	EditedAt   *time.Time   `db:"editedAt,omitempty"   json:"editedAt,omitempty"`
	Deleted    bool         `db:"deleted"              json:"deleted"`
	ReadBy     []string     `db:"readBy,omitempty"     json:"readBy,omitempty"`
	Status     string       `db:"status,omitempty"     json:"status,omitempty"` // e.g. "sent", "read"
	Nonce      string       `db:"nonce,omitempty"      json:"nonce,omitempty"`
	Seq        int64        `db:"seq,omitempty"        json:"seq,omitempty"`
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
