// File: internal/pay/refundsModels.go

package pay

import "time"

// RefundRequest represents a user's request to refund an order
type OrderRefundRequest struct {
	ID            string `db:"_id,omitempty" json:"id"`
	OrderID       string `db:"order_id" json:"order_id"`                                 // Order being refunded
	UserID        string `db:"userid" json:"userid"`                                     // User requesting refund
	OrderType     string `db:"order_type" json:"order_type"`                             // "regular" or "farm"
	Amount        int64  `db:"amount" json:"amount"`                                     // Refund amount in paise
	Reason        string `db:"reason" json:"reason"`                                     // Reason for refund request
	Status        string `db:"status" json:"status"`                                     // "pending", "approved", "rejected", "completed"
	TransactionID string `db:"transaction_id,omitempty" json:"transaction_id,omitempty"` // Created refund transaction ID

	// Admin review info
	ReviewedBy  string    `db:"reviewed_by,omitempty" json:"reviewed_by,omitempty"`   // Admin user ID who reviewed
	ReviewedAt  time.Time `db:"reviewed_at,omitempty" json:"reviewed_at,omitempty"`   // When refund was reviewed
	ReviewNotes string    `db:"review_notes,omitempty" json:"review_notes,omitempty"` // Admin notes on refund

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// RefundRequestFilter helps filter refund requests
type RefundRequestFilter struct {
	UserID    string
	OrderID   string
	Status    string
	OrderType string
	Skip      int
	Limit     int
}
