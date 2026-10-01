// File: internal/cart/cartModels.go

package cart

import (
	"time"
)

// // CartItem represents a single item in the user's cart.
//
//	type CartItem struct {
//		CartItemID string    `json:"cartItemId" db:"_id,omitempty"`
//		UserID     string    `json:"userid" db:"userid"`
//		Category   string    `json:"category" db:"category"`
//		ItemID     string    `json:"itemId" db:"itemId"`
//		ItemName   string    `json:"itemName" db:"itemName"`
//		ItemType   string    `json:"itemType,omitempty" db:"itemType,omitempty"`
//		Unit       string    `json:"unit,omitempty" db:"unit,omitempty"`
//		Discount   int64     `json:"discount,omitempty" db:"discount,omitempty"`
//		EntityID   string    `json:"entityId,omitempty" db:"entityId,omitempty"`
//		EntityName string    `json:"entityName,omitempty" db:"entityName,omitempty"`
//		EntityType string    `json:"entityType,omitempty" db:"entityType,omitempty"`
//		Quantity   int       `json:"quantity" db:"quantity"`
//		Price      int64     `json:"price,omitempty" db:"price,omitempty"` // CRITICAL FIX: Changed from float64 to int64 (stored in paise)
//		AddedAt    time.Time `json:"addedAt" db:"addedAt"`
//	}
type CartItem struct {
	ID string `db:"_id,omitempty" json:"id,omitempty"`

	UserID string `db:"userid" json:"-"`

	ItemID   string `db:"itemId" json:"itemId"`
	ItemType string `db:"itemType" json:"itemType"`

	EntityID   string `db:"entityId,omitempty" json:"entityId,omitempty"`
	EntityType string `db:"entityType,omitempty" json:"entityType,omitempty"`

	ItemName string `db:"itemName" json:"itemName"`

	Quantity int `db:"quantity" json:"quantity"`

	/*
		Price and Discount are stored as integer minor units.

		For example:
		₹199.50 -> 19950
	*/
	Price    int64 `db:"price" json:"price"`
	Discount int64 `db:"discount" json:"discount"`

	Unit     string `db:"unit,omitempty" json:"unit,omitempty"`
	Category string `db:"category,omitempty" json:"category,omitempty"`

	AddedAt   time.Time `db:"addedAt" json:"addedAt"`
	UpdatedAt time.Time `db:"updatedAt" json:"updatedAt"`
}

// CheckoutSession represents a pre-order session, grouped by category.
type CheckoutSession struct {
	UserID         string                `json:"userid" db:"userid"`
	Items          map[string][]CartItem `json:"items" db:"items"`
	Address        string                `json:"address" db:"address"`
	Total          int64                 `json:"total" db:"total"`       // CRITICAL FIX: Changed from float64 to int64 (stored in paise)
	Subtotal       int64                 `json:"subtotal" db:"subtotal"` // CRITICAL FIX: Changed from float64 to int64 (stored in paise)
	Tax            int64                 `json:"tax" db:"tax"`           // CRITICAL FIX: Changed from float64 to int64 (stored in paise)
	Delivery       int64                 `json:"delivery" db:"delivery"` // CRITICAL FIX: Changed from float64 to int64 (stored in paise)
	Discount       int64                 `json:"discount" db:"discount"` // CRITICAL FIX: Changed from float64 to int64 (stored in paise)
	PaymentMethod  string                `json:"paymentMethod" db:"paymentMethod"`
	PaymentDetails interface{}           `json:"paymentDetails" db:"paymentDetails"`
	CreatedAt      time.Time             `json:"createdAt" db:"createdAt"`
}

// Order represents a finalized order.
type Order struct {
	OrderID       string                `json:"orderId" db:"orderId"`
	OrderType     string                `json:"orderType" db:"orderType"`
	UserID        string                `json:"userid" db:"userid"`
	Items         map[string][]CartItem `json:"items" db:"items"` // grouped by category
	Address       string                `json:"address" db:"address"`
	PaymentMethod string                `json:"paymentMethod" db:"paymentMethod"`
	Status        string                `json:"status" db:"status"` // e.g. "pending", "completed"
	ApprovedBy    []string              `json:"approvedBy" db:"approvedBy"`
	CreatedAt     time.Time             `json:"createdAt" db:"createdAt"`
	Subtotal      int64                 `json:"subtotal" db:"subtotal"`
	Discount      int64                 `json:"discount" db:"discount"`
	Tax           int64                 `json:"tax" db:"tax"`
	Delivery      int64                 `json:"delivery" db:"delivery"`
	Total         int64                 `json:"total" db:"total"`
	Name          string                `json:"name" db:"name"`
	Phone         string                `json:"phone" db:"phone"`
}

type FarmOrder struct {
	OrderID         string                `db:"orderid,omitempty"  json:"orderid"`
	UserID          string                `db:"userid"         json:"userid"`
	FarmID          string                `db:"farmid"         json:"farmid"`
	CropID          string                `db:"cropid"         json:"cropid"`
	Quantity        int                   `db:"quantity"       json:"quantity"`
	PriceAtPurchase float64               `db:"priceAtPurchase" json:"priceAtPurchase"`
	CreatedAt       time.Time             `db:"createdAt"       json:"createdAt"`
	Status          OrderStatus           `db:"status"       json:"status"`
	ApprovedBy      []string              `db:"approved"       json:"approved"`
	Items           map[string][]CartItem `json:"items" db:"items"`
	Subtotal        int64                 `json:"subtotal" db:"subtotal"`
	Discount        int64                 `json:"discount" db:"discount"`
	Tax             int64                 `json:"tax" db:"tax"`
	Delivery        int64                 `json:"delivery" db:"delivery"`
	Total           int64                 `json:"total" db:"total"`
	Address         string                `json:"address" db:"address"`
	Name            string                `json:"name" db:"name"`
	Phone           string                `json:"phone" db:"phone"`
}

type OrderStatus string

const (
	OrderActive   OrderStatus = "active"
	OrderRejected OrderStatus = "rejected"
	OrderClosed   OrderStatus = "closed"
)
