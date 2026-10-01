// File: internal/reports/reportModels.go

package reports

import (
	"time"
)

type Report struct {
	ReportID    string    `db:"reportid,omitempty" json:"id"`
	ReportedBy  string    `json:"reportedBy"  db:"reportedBy"`
	TargetID    string    `json:"targetId"    db:"targetId"`
	TargetType  string    `json:"targetType"  db:"targetType"`
	Reason      string    `json:"reason"      db:"reason"`
	Notes       string    `json:"notes,omitempty"      db:"notes,omitempty"`
	Status      string    `json:"status"      db:"status"`
	ReviewedBy  string    `json:"reviewedBy,omitempty"  db:"reviewedBy,omitempty"`
	ReviewNotes string    `json:"reviewNotes,omitempty" db:"reviewNotes,omitempty"`
	CreatedAt   time.Time `json:"createdAt"   db:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"   db:"updatedAt"`

	// New fields for parent reference
	ParentType string `json:"parentType,omitempty" db:"parentType,omitempty"`
	ParentID   string `json:"parentId,omitempty"   db:"parentId,omitempty"`

	// New field to indicate whether the reporter has been notified
	Notified bool `json:"notified" db:"notified"`
}
