// File: internal/cart/models.go

package cart

import (
	"time"
)

type removeFromCartRequest struct {
	ItemID     string `json:"itemId"`
	ItemType   string `json:"itemType"`
	Category   string `json:"category,omitempty"`
	EntityID   string `json:"entityId,omitempty"`
	EntityType string `json:"entityType,omitempty"`
}

type placeOrderRequest struct {
	Address       string                `json:"address"`
	Items         map[string][]CartItem `json:"items"`
	PaymentMethod string                `json:"paymentMethod"`
	Coupon        string                `json:"coupon"`
}

type combinedOrder struct {
	OrderID       string                `db:"orderId" json:"orderId"`
	OrderType     string                `json:"orderType"` // "regular" or "farm"
	UserID        string                `db:"userid" json:"userid"`
	FarmID        string                `json:"farmId,omitempty"`
	Items         map[string][]CartItem `db:"items" json:"items,omitempty"`
	Address       string                `db:"address" json:"address,omitempty"`
	PaymentMethod string                `db:"paymentMethod" json:"paymentMethod,omitempty"`
	Total         int64                 `db:"total" json:"total"` // In paise
	Status        string                `db:"status" json:"status"`
	CreatedAt     time.Time             `db:"createdAt" json:"createdAt"`
	ApprovedBy    []string              `db:"approvedBy" json:"approvedBy,omitempty"`
}

type Coupon struct {
	Code       string    `db:"code" json:"code"`
	Discount   float64   `db:"discount" json:"discount"` // % value
	ExpiresAt  time.Time `db:"expiresat" json:"expiresat"`
	Active     bool      `db:"active" json:"active"`
	EntityID   string    `db:"entityid" json:"entityid"`
	EntityType string    `db:"entitytype" json:"entitytype"`
}

type CouponRequest struct {
	Code       string  `json:"code"`
	Cart       float64 `json:"cart"`
	EntityID   string  `json:"entityid"`
	EntityType string  `json:"entitytype"`
}

//nolint:unused
type createSessionPayload struct {
	Address       string                `json:"address"`
	Items         map[string][]CartItem `json:"items"`
	PaymentMethod string                `json:"paymentmethod"`
	Coupon        string                `json:"coupon"`
}

// ItemDetails represents item metadata fetched across various entity collections.
type ItemDetails struct {
	Name       string  `json:"name" db:"name"`
	Type       string  `json:"type" db:"type"`
	Category   string  `json:"category" db:"category"`
	Price      float64 `json:"price" db:"price"`
	Discount   float64 `json:"discount" db:"discount"`
	Unit       string  `json:"unit" db:"unit"`
	EntityID   string  `json:"entity_id" db:"entity_id"`
	EntityName string  `json:"entity_name" db:"entity_name"`
	EntityType string  `json:"entity_type" db:"entity_type"`
	Available  int     `json:"available" db:"available"`
}
