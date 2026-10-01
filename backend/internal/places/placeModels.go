// File: internal/places/placeModels.go

package places

import (
	"scav/internal/verticals/media"
	"time"
)

type Place struct {
	PlaceID           string            `json:"placeid" db:"placeid"`
	Name              string            `json:"name" db:"name"`
	ShortDesc         string            `json:"short_desc" db:"short_desc"`
	Description       string            `json:"description" db:"description"`
	Place             string            `json:"place" db:"place"`
	Capacity          int               `json:"capacity" db:"capacity"`
	Date              time.Time         `json:"date" db:"date"`
	Address           string            `json:"address" db:"address"`
	CreatedBy         string            `json:"createdBy,omitempty" db:"createdBy,omitempty"`
	OrganizerName     string            `json:"organizer_name" db:"organizer_name"`
	OrganizerContact  string            `json:"organizer_contact" db:"organizer_contact"`
	Category          string            `json:"category" db:"category"`
	Banner            string            `json:"banner" db:"banner"`
	WebsiteURL        string            `json:"website_url" db:"website_url"`
	Status            string            `json:"status" db:"status"`
	AccessibilityInfo string            `json:"accessibility_info" db:"accessibility_info"`
	SocialMediaLinks  []string          `json:"social_links" db:"social_links"`
	Tags              []string          `json:"tags" db:"tags"`
	CustomFields      map[string]any    `json:"custom_fields" db:"custom_fields"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at"`
	City              string            `json:"city,omitempty" db:"city,omitempty"`
	Country           string            `json:"country,omitempty" db:"country,omitempty"`
	ZipCode           string            `json:"zipCode,omitempty" db:"zipCode,omitempty"`
	Jobs              string            `json:"jobs,omitempty" db:"jobs,omitempty"`
	Location          Coordinates       `json:"location" db:"location,omitempty"`
	Phone             string            `json:"phone,omitempty" db:"phone,omitempty"`
	Website           string            `json:"website,omitempty" db:"website,omitempty"`
	IsOpen            bool              `json:"isopen,omitempty" db:"isopen,omitempty"`
	Distance          float64           `json:"distance,omitempty" db:"distance,omitempty"`
	Views             int               `json:"views,omitempty" db:"views,omitempty"`
	ReviewCount       int               `json:"reviewcount,omitempty" db:"reviewcount,omitempty"`
	SocialLinks       map[string]string `json:"socialLinks,omitempty" db:"socialLinks,omitempty"`
	UpdatedBy         string            `json:"updatedBy,omitempty" db:"updatedBy,omitempty"`
	DeletedAt         *time.Time        `json:"deletedAt,omitempty" db:"deletedAt,omitempty"`
	Amenities         []string          `json:"amenities,omitempty" db:"amenities,omitempty"`
	Events            []string          `json:"events,omitempty" db:"events,omitempty"`
	OperatingHours    []string          `json:"operatinghours,omitempty" db:"operatinghours,omitempty"`
	Keywords          []string          `json:"keywords,omitempty" db:"keywords,omitempty"`
}

type PlaceStatus string

const (
	PlaceActive   PlaceStatus = "active"
	PlaceInactive PlaceStatus = "inactive"
	PlaceClosed   PlaceStatus = "closed"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude,omitempty" db:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty" db:"longitude,omitempty"`
}

type CheckIn struct {
	UserID    string        `json:"userid,omitempty" db:"userid,omitempty"`
	PlaceID   string        `json:"placeId,omitempty" db:"placeId,omitempty"`
	Timestamp time.Time     `json:"timestamp,omitempty" db:"timestamp,omitempty"`
	Comment   string        `json:"comment,omitempty" db:"comment,omitempty"`
	Rating    float64       `json:"rating,omitempty" db:"rating,omitempty"` // Optional
	Medias    []media.Media `json:"images,omitempty" db:"images,omitempty"` // Optional
}

type PlaceVersion struct {
	PlaceID   string            `json:"placeId,omitempty" db:"placeId,omitempty"`
	Version   int               `json:"version,omitempty" db:"version,omitempty"`
	Data      Place             `json:"data,omitempty" db:"data,omitempty"`
	UpdatedAt time.Time         `json:"updatedAt,omitempty" db:"updatedAt,omitempty"`
	UpdatedBy string            `json:"updatedBy,omitempty" db:"updatedBy,omitempty"`
	Changes   map[string]string `json:"changes,omitempty" db:"changes,omitempty"`
}

type OperatingHours struct {
	Day          []string `json:"day,omitempty" db:"day,omitempty"`
	OpeningHours []string `json:"opening,omitempty" db:"opening,omitempty"`
	ClosingHours []string `json:"closing,omitempty" db:"closing,omitempty"`
	TimeZone     string   `json:"timeZone,omitempty" db:"timeZone,omitempty"`
}

type Tag struct {
	ID     string   `json:"id,omitempty" db:"_id,omitempty"`
	Name   string   `json:"name,omitempty" db:"name,omitempty"`
	Places []string `json:"places,omitempty" db:"places,omitempty"` // List of Place IDs tagged with this keyword
}

const (
	PlaceStatusActive     = "active"
	PlaceStatusClosed     = "closed"
	PlaceStatusRenovation = "under renovation"
)

type PlacesResponse struct {
	PlaceID        string   `json:"placeid"`
	Name           string   `json:"name"`
	ShortDesc      string   `json:"short_desc"`
	Address        string   `json:"address,omitempty"`
	Distance       float64  `json:"distance,omitempty"`
	OperatingHours []string `json:"operatinghours,omitempty"`
	Category       string   `json:"category"`
	Tags           []string `json:"tags"`
	Banner         string   `json:"banner"`
}
