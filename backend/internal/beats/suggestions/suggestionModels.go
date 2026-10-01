// File: internal/beats/suggestions/suggestionModels.go

package suggestions

type PlaceSuggestion struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Banner   string `json:"banner"`
	Category string `json:"category"`
}

type UserSuggestion struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// UserProfileResponse defines the structure for the user profile response
type UserSuggest struct {
	Username    string `json:"username" db:"username"`
	UserID      string `json:"userid" db:"userid"`
	IsFollowing bool
	Bio         string `json:"bio,omitempty" db:"bio,omitempty"`
}

type Suggestion struct {
	ID          string `json:"id" db:"suggesstionid,omitempty"`
	Type        string `json:"type" db:"type"` // e.g., "place" or "event"
	Title       string `json:"title" db:"title"`
	Description string `json:"description,omitempty" db:"description,omitempty"`
	Name        string `json:"name"`
}
