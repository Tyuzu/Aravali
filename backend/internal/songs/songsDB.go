package songs

import (
	"context"
	"scav/config"
	"scav/infra/db"
	"scav/internal/artists"
	"time"
)

var (
	SongsCollection = config.Collections.SongsCollection
)

func ListPublishedSongsByArtist(ctx context.Context, d db.Database, artistID string, result *[]ArtistSong) error {
	return d.FindMany(ctx, SongsCollection, map[string]any{"artistid": artistID, "published": true}, result)
}

func FindSongsByArtist(ctx context.Context, d db.Database, artistID string, result *[]ArtistSong) error {
	return ListPublishedSongsByArtist(ctx, d, artistID, result)
}

func FindSongByArtistAndID(ctx context.Context, d db.Database, artistID, songID string, result *ArtistSong) error {
	return d.FindOne(ctx, SongsCollection, map[string]any{"artistid": artistID, "songid": songID}, result)
}

func SaveSong(ctx context.Context, d db.Database, song *ArtistSong) error {
	return d.Insert(ctx, SongsCollection, song)
}

func InsertArtistSong(ctx context.Context, d db.Database, song *ArtistSong) error {
	return SaveSong(ctx, d, song)
}

// UpdateArtistSongFromPayload maps optional payload fields to query maps and calls UpdateArtistSong.
func UpdateArtistSongFromPayload(ctx context.Context, d db.Database, artistID, songID string, payload songPayload) (any, error) {
	updateFields := map[string]any{}

	assignIfPresent := func(field string, val *string) {
		if val != nil {
			updateFields[field] = *val
		}
	}

	assignIfPresent("title", payload.Title)
	assignIfPresent("genre", payload.Genre)
	assignIfPresent("duration", payload.Duration)
	assignIfPresent("description", payload.Description)
	assignIfPresent("audioUrl", payload.Audio)
	assignIfPresent("poster", payload.Poster)
	assignIfPresent("audioextn", payload.AudioExtn)
	assignIfPresent("posterextn", payload.PosterExtn)

	if len(updateFields) == 0 {
		return nil, artists.ErrNoFieldsToUpdate
	}

	updateFields["updatedAt"] = time.Now()

	return UpdateArtistSong(ctx, d, artistID, songID, updateFields)
}

func UpdateArtistSong(ctx context.Context, d db.Database, artistID, songID string, update map[string]any) (any, error) {
	return d.Update(ctx, SongsCollection, map[string]any{"artistid": artistID, "songid": songID}, map[string]any{"$set": update})
}

func DeleteArtistSong(ctx context.Context, d db.Database, artistID, songID string) error {
	_, err := d.Delete(ctx, SongsCollection, map[string]any{"artistid": artistID, "songid": songID})
	return err
}
