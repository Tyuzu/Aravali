package baito

import (
	"time"
)

type Baito struct {
	BaitoId          string     `json:"baitoid"`
	EntityType       string     `json:"entityType"`
	EntityID         string     `json:"entityId"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Category         string     `json:"category"`
	SubCategory      string     `json:"subcategory"`
	Location         string     `json:"location"`
	Wage             string     `json:"wage"`
	Phone            string     `json:"phone"`
	Requirements     string     `json:"requirements"`
	Banner           string     `json:"banner,omitempty"`
	Images           []string   `json:"images"`
	WorkHours        string     `json:"workHours"`
	Benefits         string     `json:"benefits,omitempty"`
	Email            string     `json:"email,omitempty"`
	Tags             []string   `json:"tags,omitempty"`
	Duration         string     `json:"duration,omitempty"`
	LastDateToApply  *time.Time `json:"lastdate,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt,omitempty"`
	OwnerID          string     `json:"ownerId"`
	ApplicationCount int        `json:"applicationcount"`
}

type BaitoApplication struct {
	BaitoAppId  string    `json:"baitoappid"`
	BaitoID     string    `json:"baitoid"`
	UserID      string    `json:"userid"`
	Username    string    `json:"username"`
	Pitch       string    `json:"pitch"`
	SubmittedAt time.Time `json:"submittedAt"`
}

/* ---------- MODELS ---------- */

type BaitosResponse struct {
	BaitoId         string     `json:"baitoid"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Category        string     `json:"category"`
	SubCategory     string     `json:"subcategory"`
	Location        string     `json:"location"`
	Wage            string     `json:"wage"`
	Requirements    string     `json:"requirements"`
	BannerURL       string     `json:"banner,omitempty"`
	WorkHours       string     `json:"workHours"`
	Duration        string     `json:"duration,omitempty"`
	LastDateToApply *time.Time `json:"lastdate,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	OwnerID         string     `json:"ownerId"`
}
