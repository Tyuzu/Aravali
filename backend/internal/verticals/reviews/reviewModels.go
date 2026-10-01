// File: internal/verticals/reviews/reviewModels.go

package reviews

import "time"

type Review struct {
	ReviewID string `json:"reviewid" db:"reviewid"`
	UserID   string `json:"userid" db:"userid"`

	EntityType string `json:"entityType" db:"entityType"`
	EntityID   string `json:"entityId" db:"entityId"`

	Rating  int    `json:"rating" db:"rating"`
	Comment string `json:"comment" db:"comment"`

	Likes    int `json:"likes,omitempty" db:"likes,omitempty"`
	Dislikes int `json:"dislikes,omitempty" db:"dislikes,omitempty"`

	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"updatedAt"`
}
