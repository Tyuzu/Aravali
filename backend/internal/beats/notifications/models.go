// File: internal/beats/notifications/models.go

// models.go
package notifications

import "time"

type Notification struct {
	NotificationID string    `json:"id" db:"notificationid"`
	UserID         string    `json:"userid" db:"userid"`
	Title          string    `json:"title" db:"title"`
	Message        string    `json:"message" db:"message"`
	Type           string    `json:"type" db:"type"` // e.g., "system", "like", "comment"
	IsRead         bool      `json:"isRead" db:"is_read"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

type NotificationPreferences struct {
	UserID      string    `json:"userid" db:"userid"`
	EmailNotifs bool      `json:"emailNotifs" db:"email_notifs"`
	PushNotifs  bool      `json:"pushNotifs" db:"push_notifs"`
	InAppNotifs bool      `json:"inAppNotifs" db:"in_app_notifs"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
