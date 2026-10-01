// File: internal/beats/activity/activityModels.go

package activity

import "time"

type Activity struct {
	// Username     string              `json:"username,omitempty" db:"username,omitempty"`
	PlaceID      string    `json:"placeId,omitempty" db:"placeId,omitempty"`
	Action       string    `json:"action,omitempty" db:"action,omitempty"`
	PerformedBy  string    `json:"performedBy,omitempty" db:"performedBy,omitempty"`
	Timestamp    time.Time `json:"timestamp,omitempty" db:"timestamp,omitempty"`
	Details      string    `json:"details,omitempty" db:"details,omitempty"`
	IPAddress    string    `json:"ipAddress,omitempty" db:"ipAddress,omitempty"`
	DeviceInfo   string    `json:"deviceInfo,omitempty" db:"deviceInfo,omitempty"`
	ActivityID   string    `json:"activityid" db:"activityid,omitempty"`
	UserID       string    `json:"userid" db:"userid"`
	ActivityType string    `json:"activity_type" db:"activity_type"` // e.g., "follow", "review", "buy"
	EntityID     string    `json:"entity_id,omitempty" db:"entity_id,omitempty"`
	EntityType   *string   `json:"entity_type,omitempty" db:"entity_type,omitempty"` // "event", "place", or null
}
