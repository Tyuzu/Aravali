package userdata

import "time"

type UserData struct {
	UserID     string `json:"userid" bson:"userid"`
	EntityID   string `json:"entity_id" bson:"entity_id"`
	EntityType string `json:"entity_type" bson:"entity_type"`
	ItemID     string `json:"item_id" bson:"item_id"`
	ItemType   string `json:"item_type" bson:"item_type"`
	CreatedAt  string `json:"created_at" bson:"created_at"`
}

// Data Models
type postDoc struct {
	PostID    string    `bson:"postid"`
	Title     string    `bson:"title"`
	Thumb     string    `bson:"thumb"`
	CreatedBy string    `bson:"createdBy"`
	Username  string    `bson:"username"`
	CreatedAt time.Time `bson:"createdAt"`
	Blocks    []struct {
		Type    string `bson:"type"`
		URL     string `bson:"url"`
		Caption string `bson:"caption"`
	} `bson:"blocks"`
}
