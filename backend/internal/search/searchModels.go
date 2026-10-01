// File: internal/search/searchModels.go

package search

import (
	"time"
)

type MEvent struct {
	EventID     string    `json:"eventid"`
	Title       string    `json:"title"`
	Location    string    `json:"location"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Image       string    `json:"banner_image"`
}

type MPlace struct {
	PlaceID     string `json:"placeid"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Image       string `json:"banner"`
	CreatedAt   string `json:"created_at"`
}

// Result represents a single search result.
type Result struct {
	Placeid     string    `json:"placeid" db:"placeid"`
	Eventid     string    `json:"eventid" db:"eventid"`
	Businessid  string    `json:"businessid" db:"businessid"`
	Userid      string    `json:"userid" db:"userid"`
	Type        string    `json:"type" db:"type"`
	Location    string    `json:"location" db:"location"`
	Address     string    `json:"address" db:"address"`
	Category    string    `json:"category" db:"category"`
	Date        time.Time `json:"date" db:"date"`
	Price       string    `json:"price" db:"price"`
	Description string    `json:"description" db:"description"`
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Contact     string    `json:"contact,omitempty"`
	Image       string    `json:"image,omitempty"`
	Link        string    `json:"link,omitempty"`
}
