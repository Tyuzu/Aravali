// File: internal/workers/workerModels.go

package workers

type BaitoWorker struct {
	UserID        string   `json:"userid" db:"userid"`
	BaitoWorkerId string   `json:"baitoWorkerId" db:"baitoWorkerId"`
	Name          string   `json:"name" db:"name"`
	Age           int      `json:"age" db:"age"`
	Phone         string   `json:"phone" db:"phone"`
	Location      string   `json:"location" db:"location"`
	Preferred     []string `json:"preferredRoles" db:"preferredRoles"`
	Bio           string   `json:"bio" db:"bio"`
	Avatar        string   `json:"avatar" db:"avatar"`
	Email         string   `json:"email,omitempty" db:"email,omitempty"`
	Experience    string   `json:"experience,omitempty" db:"experience,omitempty"`
	Skills        string   `json:"skills,omitempty" db:"skills,omitempty"`
	Availability  string   `json:"availability,omitempty" db:"availability,omitempty"`
	ExpectedWage  string   `json:"expectedWage,omitempty" db:"expectedWage,omitempty"`
	Languages     string   `json:"languages,omitempty" db:"languages,omitempty"`
	Documents     []string `json:"documents,omitempty" db:"documents,omitempty"`
	CreatedAt     int64    `json:"createdAt" db:"createdAt"`
	UpdatedAt     int64    `json:"updatedAt,omitempty" db:"updatedAt,omitempty"`
}
type BaitoWorkersResponse struct {
	UserID        string   `json:"userid" db:"userid"`
	BaitoWorkerId string   `json:"baitoWorkerId" db:"baitoWorkerId"`
	Name          string   `json:"name" db:"name"`
	Age           int      `json:"age" db:"age"`
	Phone         string   `json:"phone" db:"phone"`
	Location      string   `json:"location" db:"location"`
	Preferred     []string `json:"preferredRoles" db:"preferredRoles"`
	Bio           string   `json:"bio" db:"bio"`
	ProfilePic    string   `json:"profilePic" db:"profilePic"`
	CreatedAt     int64    `json:"createdAt" db:"createdAt"`
}
