// File: internal/merch/merchModels.go

package merch

import "time"

type Merch struct {
	MerchID string `json:"merchid" db:"merchid"`
	// EventID     string             `json:"eventid" db:"eventid"` // Reference to Event ID
	Name        string     `json:"name" db:"name"`
	Slug        string     `json:"slug,omitempty" db:"slug,omitempty"`         // URL-friendly name (e.g. "concert-tshirt")
	SKU         string     `json:"sku,omitempty" db:"sku,omitempty"`           // Stock Keeping Unit, unique per product
	Category    string     `json:"category,omitempty" db:"category,omitempty"` // e.g. “T-Shirts”, “Accessories”
	Price       float64    `json:"price" db:"price"`
	Discount    float64    `json:"discount,omitempty" db:"discount,omitempty"`         // e.g. 0.10 for 10% off
	Stock       int        `json:"stock" db:"stock"`                                   // Number of items available
	StockStatus string     `json:"stock_status,omitempty" db:"stock_status,omitempty"` // e.g. “In Stock”, “Out of Stock”, “Preorder”
	MerchPhoto  string     `json:"merch_pic" db:"merch_pic"`
	Gallery     []string   `json:"gallery,omitempty" db:"gallery,omitempty"` // Additional image filenames
	EntityID    string     `json:"entity_id" db:"entity_id"`
	EntityType  string     `json:"entity_type" db:"entity_type"` // “event” or “place”
	Description string     `json:"description,omitempty" db:"description,omitempty"`
	ShortDesc   string     `json:"short_desc,omitempty" db:"short_desc,omitempty"` // One-line summary
	Rating      float64    `json:"rating,omitempty" db:"rating,omitempty"`         // Average rating (0.0–5.0)
	ReviewCount int        `json:"review_count,omitempty" db:"review_count,omitempty"`
	Weight      float64    `json:"weight,omitempty" db:"weight,omitempty"`         // In kilograms/pounds
	Dimensions  string     `json:"dimensions,omitempty" db:"dimensions,omitempty"` // e.g. “30×20×2 cm”
	Tags        []string   `json:"tags,omitempty" db:"tags,omitempty"`             // e.g. ["rock", "tshirt"]
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updatedAt"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deletedAt,omitempty"` // Soft delete timestamp
	UserID      string     `db:"userid" json:"userid"`
}
