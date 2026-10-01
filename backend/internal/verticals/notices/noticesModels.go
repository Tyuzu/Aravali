// File: internal/verticals/notices/noticesModels.go

package notices

import (
	"time"
)

// Correct Notice model (in case you define here)
type Notice struct {
	NoticeID   string    `db:"noticeid,omitempty" json:"noticeid"`
	EntityType string    `db:"entityType" json:"entityType"`
	EntityId   string    `db:"entityId" json:"entityId"`
	Title      string    `db:"title" json:"title"`
	Content    string    `db:"content,omitempty" json:"content,omitempty"`
	Summary    string    `db:"summary" json:"summary"`
	CreatedBy  string    `db:"createdBy" json:"createdBy"`
	CreatedAt  time.Time `db:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time `db:"updatedAt" json:"updatedAt"`
}
