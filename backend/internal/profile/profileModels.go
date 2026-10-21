package profile

// ProfileUpdateRequest defines the editable profile fields for a user.
type ProfileUpdateRequest struct {
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	Name        string `json:"name,omitempty"`
	Bio         string `json:"bio,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// PublicProfileResponse is returned for public profile reads.
type PublicProfileResponse struct {
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Bio      string `json:"bio,omitempty"`
	IsPublic bool   `json:"is_public"`
}
