// File: internal/beats/userdata/userdataModels.go

package userdata

import "time"

type UserData struct {
	UserID     string `json:"userid" db:"userid"`
	EntityID   string `json:"entity_id" db:"entity_id"`
	EntityType string `json:"entity_type" db:"entity_type"`
	ItemID     string `json:"item_id" db:"item_id"`
	ItemType   string `json:"item_type" db:"item_type"`
	CreatedAt  string `json:"created_at" db:"created_at"`
}

// Data Models
type postDoc struct {
	PostID    string    `db:"postid"`
	Title     string    `db:"title"`
	Thumb     string    `db:"thumb"`
	CreatedBy string    `db:"createdBy"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"createdAt"`
	Blocks    []struct {
		Type    string `db:"type"`
		URL     string `db:"url"`
		Caption string `db:"caption"`
	} `db:"blocks"`
}
