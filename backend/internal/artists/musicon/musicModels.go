// File: internal/artists/musicon/musicModels.go

package musicon

import "time"

// --------------------------- Structs ---------------------------

type Album struct {
	ReleaseDate string   `json:"releaseDate" db:"releaseDate"`
	Description string   `json:"description" db:"description"`
	Published   bool     `json:"published" db:"published"`
	Title       string   `json:"title" db:"title"`
	ArtistID    string   `json:"artistid" db:"artistid"`
	AlbumID     string   `json:"albumid" db:"albumid"`
	Songs       []string `json:"songs" db:"songs"`
	CoverURL    string   `json:"coverUrl,omitempty" db:"coverUrl,omitempty"`
}

type Playlist struct {
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	UserID        string    `json:"userid" db:"userid"`
	PlaylistID    string    `json:"playlistid" db:"playlistid"`
	Songs         []string  `json:"songs" db:"songs"`
	CreatedAt     time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updatedAt"`
	Duration      int       `json:"duration" db:"duration"`
	IsCompilation bool      `json:"isCompilation" db:"isCompilation"`
	Copyrights    string    `json:"copyrights" db:"copyrights"`
	CoverURL      string    `db:"coverUrl,omitempty"`
}

type Song struct {
	SongID      string    `json:"songid" db:"songid,omitempty"`
	ArtistID    string    `json:"artistid" db:"artistid,omitempty"`
	Title       string    `json:"title" db:"title"`
	Genre       string    `json:"genre" db:"genre"`
	Duration    string    `json:"duration" db:"duration"`
	Description string    `json:"description,omitempty" db:"description,omitempty"`
	AudioURL    string    `json:"audioUrl,omitempty" db:"audioUrl,omitempty"`
	Published   bool      `json:"published" db:"published"`
	Plays       int       `json:"plays,omitempty" db:"plays,omitempty"`
	UploadedAt  time.Time `json:"uploadedAt" db:"uploadedAt"`
	Poster      string    `db:"poster,omitempty" json:"poster,omitempty"`
	Language    string    `json:"language" db:"language"`
	AudioExtn   string    `json:"audioextn" db:"audioextn"`
	PosterExtn  string    `json:"posterextn" db:"posterextn"`
}
