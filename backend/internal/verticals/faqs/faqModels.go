// File: internal/verticals/faqs/faqModels.go

package faqs

import "time"

type FAQ struct {
	FAQID      string    `json:"faqid" db:"faqid,omitempty"`
	EntityType string    `json:"entityType" db:"entity_type"`
	EntityID   string    `json:"entityId" db:"entity_id"`
	Content    string    `json:"content" db:"content"`
	CreatedBy  string    `json:"createdBy" db:"created_by"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
	Likes      int       `json:"likes" db:"likes"`
}
