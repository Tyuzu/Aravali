// File: internal/artists/songs/songModels.go

package songs

import "time"

type ArtistSong struct {
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
