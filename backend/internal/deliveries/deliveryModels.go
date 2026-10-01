// File: internal/deliveries/deliveryModels.go

package deliveries

import "time"

type Location struct {
	Address     string    `json:"address" db:"address"`
	Lat         float64   `json:"lat" db:"lat"`
	Lng         float64   `json:"lng" db:"lng"`
	Type        string    `json:"type,omitempty" db:"type,omitempty"`               // e.g., "Point"
	Coordinates []float64 `json:"coordinates,omitempty" db:"coordinates,omitempty"` // [lng, lat]
}

type StatusHistoryItem struct {
	Status    string    `json:"status" db:"status"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	UpdatedBy string    `json:"updated_by" db:"updated_by"`
}

type Proof struct {
	ProofID   string    `json:"proofid" db:"id"`
	Type      string    `json:"type" db:"type"` // e.g. "PHOTO", "SIGNATURE"
	URL       string    `json:"url" db:"url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Delivery struct {
	DeliveryID          string              `json:"deliveryid" db:"id"`
	TenantID            string              `json:"tenantid" db:"tenantid"`
	UserID              string              `json:"userid" db:"userid"`
	DriverID            *string             `json:"driverid" db:"driverid"`
	Status              string              `json:"status" db:"status"`
	StatusHistory       []StatusHistoryItem `json:"status_history,omitempty" db:"status_history,omitempty"`
	PickupLoc           Location            `json:"pickup_loc" db:"pickup_loc"`
	DropoffLoc          Location            `json:"dropoff_loc" db:"dropoff_loc"`
	CurrentLocation     *Location           `json:"current_location,omitempty" db:"current_location,omitempty"`
	Proofs              []Proof             `json:"proofs,omitempty" db:"proofs,omitempty"`
	PublicTrackingToken string              `json:"public_tracking_token,omitempty" db:"public_tracking_token"`
	EstimatedArrival    *time.Time          `json:"estimated_arrival,omitempty" db:"estimated_arrival,omitempty"`
	CreatedAt           time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at" db:"updated_at"`
}

type GPSData struct {
	Lat       float64   `json:"lat" db:"lat"`
	Lng       float64   `json:"lng" db:"lng"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

type Driver struct {
	DriverID     string    `json:"driverid" db:"id"`
	TenantID     string    `json:"tenantid" db:"tenantid"`
	Name         string    `json:"name" db:"name"`
	IsOnline     bool      `json:"is_online" db:"is_online"`
	CurrentState string    `json:"current_state" db:"current_state"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Webhook struct {
	WebhookID string    `json:"webhookid" db:"id"`
	TenantID  string    `json:"tenantid" db:"tenantid"`
	URL       string    `json:"url" db:"url"`
	Events    []string  `json:"events" db:"events"`
	Secret    string    `json:"secret" db:"secret"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
