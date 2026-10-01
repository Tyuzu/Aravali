// File: internal/events/eventModels.go

package events

import (
	"scav/internal/vendors"
	"time"
)

type Event struct {
	EventID          string      `json:"eventid" db:"eventid"`
	Title            string      `json:"title" db:"title"`
	Description      string      `json:"description" db:"description"`
	Date             time.Time   `json:"date" db:"date"`
	PlaceID          string      `json:"placeid" db:"placeid"`
	PlaceName        string      `json:"placename" db:"placename"`
	Location         string      `json:"location" db:"location"`
	Coords           Coordinates `json:"coords" db:"coords"`
	CreatorID        string      `json:"creatorid" db:"creatorid"`
	StartDateTime    time.Time   `json:"start_date_time" db:"start_date_time"`
	EndDateTime      time.Time   `json:"end_date_time" db:"end_date_time"`
	Category         string      `json:"category" db:"category"`
	Banner           string      `json:"banner" db:"banner"`
	SeatingPlanImage string      `json:"seating" db:"seating"`
	WebsiteURL       string      `json:"website_url" db:"website_url"`
	Status           string      `json:"status" db:"status"`
	Tags             []string    `json:"tags" db:"tags"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`
	OrganizerName    string      `json:"organizer_name" db:"organizer_name"`
	OrganizerContact string      `json:"organizer_contact" db:"organizer_contact"`
	Artists          []string    `json:"artists,omitempty" db:"artists,omitempty"`
	Published        string      `json:"published,omitempty" db:"published,omitempty"`
	External         bool        `json:"external" db:"external"`
	ExternalLink     string      `json:"externallink" db:"externallink"`
	// New fields for alignment (CRITICAL FIX)
	ContactInfo  *EventContactInfo      `json:"contactInfo" db:"contact_info"`
	News         []NewsItem             `json:"news" db:"news"`
	Polls        []Poll                 `json:"polls" db:"polls"`
	LostFound    []LostFoundItem        `json:"lostfound" db:"lost_found"`
	HiredVendors []vendors.VendorHiring `json:"hired_vendors,omitempty" db:"hired_vendors,omitempty"`
	// Computed fields for frontend filters
	Prices   []float64 `json:"prices,omitempty" db:"-"`
	Currency string    `json:"currency,omitempty" db:"-"`
}

// FAQ represents a single FAQ structure
type FAQ struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// EventContactInfo represents event contact information (renamed to avoid conflicts with Farm.ContactInfo)
type EventContactInfo struct {
	Email         string `json:"email" db:"email"`
	Phone         string `json:"phone" db:"phone"`
	OrganizerName string `json:"organizer_name" db:"organizer_name"`
}

// NewsItem represents a single news update for an event
type NewsItem struct {
	ID        string    `json:"id" db:"_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

// PollOption represents a single poll option with vote count
type PollOption struct {
	Text  string `json:"text" db:"text"`
	Votes int    `json:"votes" db:"votes"`
}

// Poll represents a poll for an event
type Poll struct {
	ID       string       `json:"id" db:"_id"`
	Question string       `json:"question" db:"question"`
	Options  []PollOption `json:"options" db:"options"`
}

// LostFoundItem represents a lost or found item at an event
type LostFoundItem struct {
	ID          string `json:"id" db:"_id"`
	Type        string `json:"type" db:"type"` // "lost" or "found"
	Description string `json:"description" db:"description"`
	Contact     string `json:"contact" db:"contact"`
}

type SocialMediaLinks struct {
	Title string `json:"title"`
	Url   string `json:"Url"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude,omitempty" db:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty" db:"longitude,omitempty"`
}
