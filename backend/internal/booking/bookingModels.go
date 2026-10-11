// File: internal/booking/bookingModels.go

package booking

// Tier defines a pricing/capacity tier for bookings
// type Tier struct {
// 	ID         string   `json:"id" db:"id"`
// 	EntityType string   `json:"entityType" db:"entityType"`
// 	EntityId   string   `json:"entityId" db:"entityId"`
// 	Name       string   `json:"name" db:"name"`
// 	Price      float64  `json:"price" db:"price"`
// 	Capacity   int      `json:"capacity" db:"capacity"`
// 	TimeRange  []string `json:"timeRange,omitempty" db:"timeRange,omitempty"`   // ["09:00", "17:00"]
// 	DaysOfWeek []int    `json:"daysOfWeek,omitempty" db:"daysOfWeek,omitempty"` // 0=Sun..6=Sat
// 	Features   []string `json:"features,omitempty" db:"features,omitempty"`
// 	CreatedAt  int64    `json:"createdAt" db:"createdAt"`
// }

// // Slot represents an available time slot for booking
// type Slot struct {
// 	ID         string `json:"id" db:"id"`
// 	EntityType string `json:"entityType" db:"entityType"`
// 	EntityId   string `json:"entityId" db:"entityId"`
// 	Date       string `json:"date" db:"date"`
// 	Start      string `json:"start" db:"start"`
// 	End        string `json:"end,omitempty" db:"end,omitempty"`
// 	Capacity   int    `json:"capacity" db:"capacity"`
// 	TierId     string `json:"tierId,omitempty" db:"tierId,omitempty"`
// 	TierName   string `json:"tierName,omitempty" db:"tierName,omitempty"`
// 	CreatedAt  int64  `json:"createdAt" db:"createdAt"`
// }

// // Booking represents a user's booking of a slot or tier
// type Booking struct {
// 	ID         string  `json:"id" db:"id"`
// 	SlotId     string  `json:"slotId,omitempty" db:"slotId,omitempty"`
// 	TierId     string  `json:"tierId,omitempty" db:"tierId,omitempty"`
// 	TierName   string  `json:"tierName,omitempty" db:"tierName,omitempty"`
// 	PricePaid  float64 `json:"pricePaid,omitempty" db:"pricePaid,omitempty"`
// 	EntityType string  `json:"entityType" db:"entityType"`
// 	EntityId   string  `json:"entityId" db:"entityId"`
// 	UserId     string  `json:"userid" db:"userid"`
// 	Date       string  `json:"date" db:"date"`
// 	Start      string  `json:"start" db:"start"`
// 	End        string  `json:"end,omitempty" db:"end,omitempty"`
// 	Status     string  `json:"status" db:"status"` // pending, confirmed, cancelled
// 	CreatedAt  int64   `json:"createdAt" db:"createdAt"`
// }

// // DateCap represents the capacity limit for a specific date
// type DateCap struct {
// 	EntityType string `json:"entityType" db:"entityType"`
// 	EntityId   string `json:"entityId" db:"entityId"`
// 	Date       string `json:"date" db:"date"`
// 	Capacity   int    `json:"capacity" db:"capacity"`
// }

// ---------- Models ----------
type Slot struct {
	ID         string `json:"id" db:"id"`
	EntityType string `json:"entityType" db:"entityType"`
	EntityId   string `json:"entityId" db:"entityId"`
	Date       string `json:"date" db:"date"`
	Start      string `json:"start" db:"start"`
	End        string `json:"end,omitempty" db:"end,omitempty"`
	Capacity   int    `json:"capacity" db:"capacity"`
	TierId     string `json:"tierId,omitempty" db:"tierId,omitempty"`
	TierName   string `json:"tierName,omitempty" db:"tierName,omitempty"`
	CreatedAt  int64  `json:"createdAt" db:"createdAt"`
}

type Booking struct {
	ID         string  `json:"id" db:"id"`
	SlotId     string  `json:"slotId,omitempty" db:"slotId,omitempty"`
	TierId     string  `json:"tierId,omitempty" db:"tierId,omitempty"`
	TierName   string  `json:"tierName,omitempty" db:"tierName,omitempty"`
	PricePaid  float64 `json:"pricePaid,omitempty" db:"pricePaid,omitempty"`
	EntityType string  `json:"entityType" db:"entityType"`
	EntityId   string  `json:"entityId" db:"entityId"`
	UserId     string  `json:"userid" db:"userid"`
	Date       string  `json:"date" db:"date"`
	Start      string  `json:"start" db:"start"`
	End        string  `json:"end,omitempty" db:"end,omitempty"`
	Status     string  `json:"status" db:"status"` // pending, confirmed, cancelled
	Seats      int     `json:"seats,omitempty" db:"seats,omitempty"`
	CreatedAt  int64   `json:"createdAt" db:"createdAt"`
}

type DateCap struct {
	EntityType string `json:"entityType" db:"entityType"`
	EntityId   string `json:"entityId" db:"entityId"`
	Date       string `json:"date" db:"date"`
	Capacity   int    `json:"capacity" db:"capacity"`
}

type Tier struct {
	ID         string   `json:"id" db:"id"`
	EntityType string   `json:"entityType" db:"entityType"`
	EntityId   string   `json:"entityId" db:"entityId"`
	Name       string   `json:"name" db:"name"`
	Price      float64  `json:"price" db:"price"`
	Capacity   int      `json:"capacity" db:"capacity"`
	TimeRange  []string `json:"timeRange,omitempty" db:"timeRange,omitempty"`   // ["09:00", "17:00"]
	DaysOfWeek []int    `json:"daysOfWeek,omitempty" db:"daysOfWeek,omitempty"` // 0=Sun..6=Sat
	Features   []string `json:"features,omitempty" db:"features,omitempty"`
	CreatedAt  int64    `json:"createdAt" db:"createdAt"`
}
