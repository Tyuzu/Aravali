// File: internal/profile/profileModels.go

package profile

import "time"

// UserProfileResponse defines the structure for the user profile response
type UserProfileResponse struct {
	UserID         string            `json:"userid" db:"userid"`
	Username       string            `json:"username" db:"username"`
	Name           string            `json:"name" db:"name"`
	Email          string            `json:"email" db:"email"`
	Bio            string            `json:"bio,omitempty" db:"bio,omitempty"`
	PhoneNumber    string            `json:"phone_number,omitempty" db:"phone_number,omitempty"`
	Avatar         string            `json:"avatar" db:"avatar"`
	Banner         string            `json:"banner" db:"banner"`
	IsFollowing    bool              `json:"is_following" db:"is_following"` // Added here
	FollowersCount int               `json:"followerscount" db:"followerscount"`
	FollowingCount int               `json:"followscount" db:"followscount"`
	SocialLinks    map[string]string `json:"social_links,omitempty" db:"social_links,omitempty"`
	Online         bool              `json:"online,omitempty"`
	LastLogin      time.Time         `json:"last_login" db:"last_login"`
}
