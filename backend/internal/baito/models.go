package baito

type BaitoRequest struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Category        string   `json:"category"`
	SubCategory     string   `json:"sub_category"`
	Location        string   `json:"location"`
	Wage            string   `json:"wage"`
	Phone           string   `json:"phone"`
	Requirements    string   `json:"requirements"`
	WorkHours       string   `json:"work_hours"`
	Benefits        string   `json:"benefits"`
	Email           string   `json:"email"`
	Tags            []string `json:"tags"`
	Duration        string   `json:"duration"`
	LastDateToApply string   `json:"last_date_to_apply"`
}

type UpdateBaitoResponse struct {
	Message string `json:"message"`
	BaitoID string `json:"baitoid"`
}
