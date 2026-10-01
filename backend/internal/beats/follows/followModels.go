// File: internal/beats/follows/followModels.go

package follows

type UserFollow struct {
	UserID    string   `json:"userid" db:"userid"`
	Follows   []string `json:"follows,omitempty" db:"follows,omitempty"`
	Followers []string `json:"followers,omitempty" db:"followers,omitempty"`
}
