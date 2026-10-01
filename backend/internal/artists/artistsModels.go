// File: internal/artists/artistsModels.go

package artists

import (
	"time"
)

type Artist struct {
	ArtistID  string            `db:"artistid,omitempty" json:"artistid"`
	Category  string            `db:"category" json:"category"`
	Name      string            `db:"name" json:"name"`
	Place     string            `db:"place" json:"place"`
	Country   string            `db:"country" json:"country"`
	Bio       string            `db:"bio" json:"bio"`
	DOB       string            `db:"dob" json:"dob"`
	Photo     string            `db:"photo" json:"photo"`
	Banner    string            `db:"banner" json:"banner"`
	Genres    []string          `db:"genres" json:"genres"`
	Socials   map[string]string `db:"socials" json:"socials"`
	EventIDs  []string          `db:"events" json:"events"`
	Members   []BandMember      `db:"members,omitempty" json:"members,omitempty"` // ✅ ADD THIS
	CreatedAt time.Time         `json:"createdAt" db:"createdAt"`
	CreatorID string            `db:"creatorid" json:"creatorid"`
}

type BandMember struct {
	MemberID        string `db:"memberid,omitempty" json:"memberid,omitempty"`
	ReferenceArtist string `db:"ref_artistid,omitempty" json:"ref_artistid,omitempty"`
	Name            string `db:"name" json:"name"`
	Role            string `db:"role,omitempty" json:"role,omitempty"`
	DOB             string `db:"dob,omitempty" json:"dob,omitempty"`
	Image           string `db:"image,omitempty" json:"image,omitempty"`
}

// ArtistEvent Struct
type ArtistEvent struct {
	EventID   string `db:"eventid,omitempty" json:"eventid"`
	ArtistID  string `db:"artistid" json:"artistid"`
	Title     string `db:"title" json:"title"`
	Date      string `db:"date" json:"date"`
	Venue     string `db:"venue" json:"venue"`
	City      string `db:"city" json:"city"`
	Country   string `db:"country" json:"country"`
	CreatorID string `db:"creatorid" json:"creatorid"`
	TicketURL string `db:"ticket_url,omitempty" json:"ticketUrl,omitempty"`
}

type ArtistAlbum struct {
	Title       string `json:"title"`
	ReleaseDate string `json:"releaseDate"`
	Description string `json:"description"`
	Published   bool   `json:"published"`
}

type ArtistPost struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	Published bool   `json:"published"`
}

type ArtistMerchItem struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Image       string  `json:"image,omitempty"`
	Visible     bool    `json:"visible"`
	MerchID     string  `json:"merchid" db:"merchid"`
}

// CreateArtistEventRequest defines the shape of the body to create an event.
type CreateArtistEventRequest struct {
	Title string `json:"title"`
	Date  string `json:"date"` // Expects "YYYY-MM-DD"
	Venue string `json:"venue"`
}

// CreateArtistEventResponse returning data after creation success.
type CreateArtistEventResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

// AddArtistToEventRequest captures standard event mapping requirements.
type AddArtistToEventRequest struct {
	EventID string `json:"eventid"`
}

// GenericMessageResponse is reused across successful operations.
type GenericMessageResponse struct {
	Message string `json:"message"`
}

type ArtistToEventRequestPayload struct {
	EventID  string `json:"eventid"`
	ArtistID string `json:"artistid"`
}

type ArtistByIDResponse struct {
	Artist
	IsSubscribed bool `json:"issubscribed"`
}
