// File: internal/vendors/vendorModels.go

package vendors

import (
	"errors"
	"time"
)

var (
	ErrVendorNotFound      = errors.New("vendor not found")
	ErrVendorAlreadyExists = errors.New("vendor already exists")
	ErrVendorAlreadyHired  = errors.New("vendor already hired for this event")
	ErrVendorNotInEvent    = errors.New("vendor not in event")
)

// Vendor represents a vendor who can be hired for events
type Vendor struct {
	VendorID     string    `json:"vendorid" db:"vendorid"`
	UserID       string    `json:"userid" db:"userid"`
	Name         string    `json:"name" db:"name"`
	Category     string    `json:"category" db:"category"`
	Description  string    `json:"description,omitempty" db:"description,omitempty"`
	Email        string    `json:"email,omitempty" db:"email,omitempty"`
	Phone        string    `json:"phone,omitempty" db:"phone,omitempty"`
	Location     string    `json:"location,omitempty" db:"location,omitempty"`
	Rating       float64   `json:"rating,omitempty" db:"rating,omitempty"`
	RatingCount  int       `json:"rating_count,omitempty" db:"rating_count,omitempty"`
	ProfileImage string    `json:"profile_image,omitempty" db:"profile_image,omitempty"`
	Portfolio    []string  `json:"portfolio,omitempty" db:"portfolio,omitempty"`
	Verified     bool      `json:"verified" db:"verified"`
	Available    bool      `json:"available" db:"available"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

// VendorHiring represents the relationship between an event and hired vendors
type VendorHiring struct {
	HiringID       string    `json:"hiringid" db:"hiringid"`
	EventID        string    `json:"eventid" db:"eventid"`
	VendorID       string    `json:"vendorid" db:"vendorid"`
	VendorName     string    `json:"vendor_name" db:"vendor_name"`
	VendorCategory string    `json:"vendor_category" db:"vendor_category"`
	HiredAt        time.Time `json:"hired_at" db:"hired_at"`
	HiredBy        string    `json:"hired_by" db:"hired_by"` // UserID of event creator/organizer
	Status         string    `json:"status" db:"status"`     // "hired", "accepted", "rejected", "completed"
	Notes          string    `json:"notes,omitempty" db:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

// VendorResponse is the response structure for vendor data
type VendorResponse struct {
	VendorID     string    `json:"vendorid"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	Description  string    `json:"description,omitempty"`
	Email        string    `json:"email,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	Location     string    `json:"location,omitempty"`
	Rating       float64   `json:"rating,omitempty"`
	RatingCount  int       `json:"rating_count,omitempty"`
	ProfileImage string    `json:"profile_image,omitempty"`
	Portfolio    []string  `json:"portfolio,omitempty"`
	Verified     bool      `json:"verified"`
	Status       string    `json:"status,omitempty" db:"status,omitempty"`
	HiringID     string    `json:"hiringid,omitempty" db:"hiringid,omitempty"`
	HiredAt      time.Time `json:"hired_at,omitempty" db:"hired_at,omitempty"`
}

// AvailabilitySlot represents a vendor's unavailable or available date range
type AvailabilitySlot struct {
	SlotID         string    `json:"slotid" db:"slotid"`
	VendorID       string    `json:"vendorid" db:"vendorid"`
	StartDate      string    `json:"start_date" db:"start_date"` // YYYY-MM-DD
	EndDate        string    `json:"end_date" db:"end_date"`     // YYYY-MM-DD
	Recurring      bool      `json:"recurring,omitempty" db:"recurring,omitempty"`
	RecurrenceRule string    `json:"recurrence_rule,omitempty" db:"recurrence_rule,omitempty"` // e.g. RFC5545 or simple rule
	Notes          string    `json:"notes,omitempty" db:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}
