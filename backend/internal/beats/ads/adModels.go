// File: internal/beats/ads/adModels.go

package ads

import "time"

type AdType string

const (
	TypeExternal AdType = "external"
	TypePost     AdType = "internal_post"
)

type Ad struct {
	ID          string    `json:"id,omitempty" db:"_id,omitempty"`
	Type        AdType    `json:"type" db:"type"`                         // "external" or "internal_post"
	PostID      string    `json:"postId,omitempty" db:"postId,omitempty"` // Reference to internal post if Type == "internal_post"
	Title       string    `json:"title,omitempty" db:"title,omitempty"`
	Description string    `json:"description,omitempty" db:"description,omitempty"`
	Image       string    `json:"image,omitempty" db:"image,omitempty"`
	Link        string    `json:"link,omitempty" db:"link,omitempty"`
	Category    string    `json:"category,omitempty" db:"category,omitempty"`
	Page        string    `json:"page,omitempty" db:"page,omitempty"`
	Position    string    `json:"position,omitempty" db:"position,omitempty"`
	Status      string    `json:"status" db:"status"` // "active", "inactive"
	CreatedAt   time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updatedAt"`
}
