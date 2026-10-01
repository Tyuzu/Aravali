// File: internal/verticals/tickets/ticketModels.go

package tickets

import (
	"time"
)

type RefundRequest struct {
	EventID     string     `db:"eventid" json:"eventID"`
	TicketID    string     `db:"ticketid" json:"ticketID"`
	UserID      string     `db:"userid" json:"userID"`
	UniqueCode  string     `db:"uniquecode" json:"uniqueCode"`
	RequestDate time.Time  `db:"requestdate" json:"requestDate"`
	Status      string     `db:"status" json:"status"` // pending, approved, rejected, refunded
	Amount      int        `db:"amount" json:"amount,omitempty"`
	ProcessedAt *time.Time `db:"processedat,omitempty" json:"processedAt,omitempty"`
	RefundedAt  *time.Time `db:"refundedat,omitempty" json:"refundedAt,omitempty"`
}

type Ticket struct {
	TicketID    string    `json:"ticketid" db:"ticketid"`
	EventID     string    `json:"eventid" db:"eventid"`
	Name        string    `json:"name" db:"name"`
	Price       int64     `json:"price" db:"price"` // CRITICAL FIX: Changed from float64 to int64 (stored in paise)
	Currency    string    `json:"currency" db:"currency"`
	Color       string    `json:"color" db:"color"`
	Quantity    int       `json:"quantity" db:"quantity"`
	EntityID    string    `json:"entity_id" db:"entity_id"`
	EntityType  string    `json:"entity_type" db:"entity_type"` // "event" or "place"
	Available   int       `json:"available" db:"available"`
	Total       int       `json:"total" db:"total"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Description string    `db:"description,omitempty" json:"description"`
	Sold        int       `db:"sold" json:"sold"`
	SeatStart   int       `db:"seatstart" json:"seatstart"`
	SeatEnd     int       `db:"seatend" json:"seatend"`
	Seats       []Seat    `db:"seats" json:"seats"` // 👈 new field
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

type Seat struct {
	SeatID     string `json:"id" db:"_id,omitempty"`
	EntityID   string `json:"entity_id" db:"entity_id"`
	EntityType string `json:"entity_type" db:"entity_type"` // e.g., "event" or "place"
	SeatNumber string `json:"seat_number" db:"seat_number"`
	UserID     string `json:"userid" db:"userid,omitempty"`
	Status     string `json:"status" db:"status"` // e.g., "booked", "available"
}

type PurchasedTicket struct {
	EventID      string    `db:"eventid" json:"eventid"`
	TicketID     string    `db:"ticketid" json:"ticketid"`
	UserID       string    `db:"userid" json:"userid"`
	BuyerName    string    `db:"buyername" json:"buyerName"`
	UniqueCode   string    `db:"uniquecode" json:"uniqueCode"`
	PurchaseDate time.Time `db:"purchasedate" json:"purchaseDate"`
	Price        int       `db:"price" json:"price"`

	// Soft delete fields
	Canceled       bool       `db:"canceled" json:"canceled"`
	CanceledAt     *time.Time `db:"canceledat,omitempty" json:"canceledAt,omitempty"`
	CanceledReason string     `db:"cancelledreason,omitempty" json:"canceledReason,omitempty"`
	Transferred    bool       `db:"transferred" json:"transferred"`
	TransferredTo  string     `db:"transferredto,omitempty" json:"transferredTo,omitempty"`
}
