// File: internal/menu/menuModels.go

package menu

import "time"

type Menu struct {
	MenuID      string    `json:"menuid" db:"menuid"`
	PlaceID     string    `json:"placeid" db:"placeid"` // Reference to Place ID
	Name        string    `json:"name" db:"name"`
	Price       float64   `json:"price" db:"price"`
	Discount    float64   `json:"discount,omitempty" db:"discount,omitempty"`
	Stock       int       `json:"stock" db:"stock"` // Number of items available
	MenuPhoto   string    `json:"menu_pic" db:"menu_pic"`
	Description string    `json:"description,omitempty" db:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UserID      string    `db:"userid" json:"userid"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}
