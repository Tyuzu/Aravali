// File: internal/verticals/media/mediaModels.go

package media

import "time"

type Media struct {
	MediaID       string    `json:"mediaid" db:"mediaid"`
	MediaGroupID  string    `json:"mediaGroupId" db:"mediaGroupId"` // new field to group multiple files
	Type          string    `json:"type" db:"type"`                 // "image", "video", "text"
	URL           string    `json:"url,omitempty" db:"url,omitempty"`
	ThumbnailURL  string    `json:"thumbnailUrl,omitempty" db:"thumbnailUrl,omitempty"`
	Caption       string    `json:"caption,omitempty" db:"caption,omitempty"`
	Description   string    `json:"description,omitempty" db:"description,omitempty"`
	CreatorID     string    `json:"creatorid" db:"creatorid"`
	LikesCount    int       `json:"likesCount" db:"likesCount"`
	CommentsCount int       `json:"commentsCount" db:"commentsCount"`
	Visibility    string    `json:"visibility,omitempty" db:"visibility,omitempty"`
	Tags          []string  `json:"tags,omitempty" db:"tags,omitempty"` // e.g., song:123, event:456
	Duration      float64   `json:"duration,omitempty" db:"duration,omitempty"`
	FileSize      int64     `json:"fileSize,omitempty" db:"fileSize,omitempty"`
	MimeType      string    `json:"mimeType,omitempty" db:"mimeType,omitempty"`
	IsFeatured    bool      `json:"isFeatured,omitempty" db:"isFeatured,omitempty"`
	EntityID      string    `json:"entityid" db:"entityid"`
	EntityType    string    `json:"entitytype" db:"entitytype"` // "event", "place", etc.
	CreatedAt     time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updatedAt"`
	UserID        string    `json:"userid" db:"userid"`
	Extn          string    `json:"extn" db:"extn"`
	CaptionLang   string    `json:"captionlang" db:"captionlang"`
}

const (
	MediaTypeImage    = "image"
	MediaTypeVideo    = "video"
	MediaTypePhoto360 = "photo360"
)
