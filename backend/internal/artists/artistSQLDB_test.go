package artists

import (
	"testing"
	"time"
)

func TestParseArtistMetadataCompatibility(t *testing.T) {
	metadata := map[string]any{
		"category":  "Singer",
		"place":     "Lagos",
		"country":   "NG",
		"dob":       "2000-01-01",
		"photo":     "https://example.com/p.jpg",
		"banner":    "https://example.com/b.jpg",
		"genres":    []any{"rock", "pop"},
		"socials":   map[string]any{"x": "@demo"},
		"creatorid": "u-123",
		"createdAt": "2024-01-02T03:04:05Z",
	}

	artist := parseArtistMetadata(metadata)
	if artist.Category != "Singer" {
		t.Fatalf("expected category Singer, got %q", artist.Category)
	}
	if artist.Place != "Lagos" {
		t.Fatalf("expected place Lagos, got %q", artist.Place)
	}
	if len(artist.Genres) != 2 {
		t.Fatalf("expected 2 genres, got %d", len(artist.Genres))
	}
	if artist.Socials["x"] != "@demo" {
		t.Fatalf("expected socials map x=@demo, got %#v", artist.Socials)
	}
	if artist.CreatorID != "u-123" {
		t.Fatalf("expected creatorid u-123, got %q", artist.CreatorID)
	}
	if artist.CreatedAt.IsZero() {
		t.Fatal("expected createdAt to be parsed")
	}
	if artist.CreatedAt.Year() != time.Date(2024, time.January, 2, 3, 4, 5, 0, utcLocation).Year() {
		t.Fatalf("unexpected createdAt value: %#v", artist.CreatedAt)
	}
}

var utcLocation = time.UTC
