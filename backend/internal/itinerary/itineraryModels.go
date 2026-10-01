// File: internal/itinerary/itineraryModels.go

package itinerary

// Itinerary represents the travel itinerary
type Itinerary struct {
	ItineraryID string  `json:"itineraryid" db:"itineraryid,omitempty"`
	UserID      string  `json:"userid" db:"userid"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	StartDate   string  `json:"start_date" db:"start_date"`
	EndDate     string  `json:"end_date" db:"end_date"`
	Status      string  `json:"status" db:"status"` // Draft/Confirmed
	Published   bool    `json:"published" db:"published"`
	ForkedFrom  *string `json:"forked_from,omitempty" db:"forked_from,omitempty"`
	Deleted     bool    `json:"-" db:"deleted,omitempty"` // Internal use only
	// the new day-by-day schedule
	Days []Day `json:"days" db:"days"`
}

// add these at the top, just below package declaration
type Visit struct {
	Location  string `json:"location" db:"location"`
	StartTime string `json:"start_time" db:"start_time"`
	EndTime   string `json:"end_time" db:"end_time"`
	// nil for the very first visit of a day
	Transport *string `json:"transport,omitempty" db:"transport,omitempty"`
}

type Day struct {
	Date   string  `json:"date" db:"date"`
	Visits []Visit `json:"visits" db:"visits"`
}
