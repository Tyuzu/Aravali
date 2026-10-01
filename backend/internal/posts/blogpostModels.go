// File: internal/posts/blogpostModels.go

package posts

import "time"

type Block struct {
	Type     string `db:"type" json:"type"`
	Content  string `db:"content,omitempty" json:"content,omitempty"`
	URL      string `db:"url,omitempty" json:"url,omitempty"`
	Alt      string `db:"alt,omitempty" json:"alt,omitempty"`
	Caption  string `db:"caption,omitempty" json:"caption,omitempty"`   // for video or image captions
	Language string `db:"language,omitempty" json:"language,omitempty"` // for code blocks
}

type BlogPost struct {
	PostID      string    `db:"postid" json:"postid"`
	Title       string    `db:"title" json:"title"`
	Category    string    `db:"category" json:"category"`
	Subcategory string    `db:"subcategory" json:"subcategory"`
	ReferenceID *string   `db:"referenceId,omitempty" json:"referenceId,omitempty"`
	Blocks      []Block   `db:"blocks" json:"blocks"`
	Thumb       string    `db:"thumb" json:"thumb"`
	CreatedBy   string    `db:"createdBy" json:"createdBy"`
	CreatedAt   time.Time `db:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `db:"updatedAt" json:"updatedAt"`
	Hashtags    []string  `db:"hashtags" json:"hashtags"`
	Type        string    `json:"type" db:"type"`
	Username    string    `json:"username" db:"username"`
}

// --- BlogPostResponse for list view ---

type BlogPostResponse struct {
	PostID      string    `db:"postid" json:"postid"`
	Title       string    `db:"title" json:"title"`
	Category    string    `db:"category" json:"category"`
	Subcategory string    `db:"subcategory" json:"subcategory"`
	ReferenceID *string   `db:"referenceId,omitempty" json:"referenceId,omitempty"`
	Thumb       string    `db:"thumb" json:"thumb"`
	CreatedBy   string    `db:"createdBy" json:"createdBy"`
	Username    string    `db:"username" json:"username"`
	CreatedAt   time.Time `db:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `db:"updatedAt" json:"updatedAt"`
}
