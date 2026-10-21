package mqevents

import "time"

const (
	SubjectUserRegistered     = "user.registered"
	SubjectUserLoggedIn       = "user.logged_in"
	SubjectUserLoggedOut      = "user.logged_out"
	SubjectUserSessionRevoked = "user.session.revoked"
	SubjectUserProfileUpdated = "user.profile.updated"
	SubjectUserProfileDeleted = "user.profile.deleted"
)

type UserRegisteredEvent struct {
	UserID    string    `json:"userid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

type UserLoggedInEvent struct {
	UserID    string    `json:"userid"`
	Username  string    `json:"username"`
	IP        string    `json:"ip,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type UserLoggedOutEvent struct {
	UserID    string    `json:"userid"`
	Username  string    `json:"username,omitempty"`
	Reason    string    `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type UserSessionRevokedEvent struct {
	UserID    string    `json:"userid"`
	Reason    string    `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type UserProfileUpdatedEvent struct {
	UserID    string    `json:"userid"`
	Username  string    `json:"username,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserProfileDeletedEvent struct {
	UserID    string    `json:"userid"`
	Username  string    `json:"username,omitempty"`
	DeletedAt time.Time `json:"deleted_at"`
}
