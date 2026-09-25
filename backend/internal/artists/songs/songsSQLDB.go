package songs

import (
	"context"
	"scav/config"
	"scav/infra/db"
	"scav/internal/artists"
	"time"
)

var (
	SongsTable = config.Tables.SongsTable
)

func SQLListPublishedSongsByArtist(ctx context.Context, d db.Database, artistID string, result *[]ArtistSong) error {
	return d.FindMany(ctx, SongsTable, map[string]any{"artistid": artistID, "published": true}, result)
}

func SQLFindSongsByArtist(ctx context.Context, d db.Database, artistID string, result *[]ArtistSong) error {
	return ListPublishedSongsByArtist(ctx, d, artistID, result)
}

func SQLFindSongByArtistAndID(ctx context.Context, d db.Database, artistID, songID string, result *ArtistSong) error {
	return d.FindOne(ctx, SongsTable, map[string]any{"artistid": artistID, "songid": songID}, result)
}

func SQLSaveSong(ctx context.Context, d db.Database, song *ArtistSong) error {
	return d.Insert(ctx, SongsTable, song)
}

func SQLInsertArtistSong(ctx context.Context, d db.Database, song *ArtistSong) error {
	return SaveSong(ctx, d, song)
}

// UpdateArtistSongFromPayload maps optional payload fields to query maps and calls UpdateArtistSong.
func SQLUpdateArtistSongFromPayload(ctx context.Context, d db.Database, artistID, songID string, payload songPayload) (any, error) {
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

func SQLUpdateArtistSong(ctx context.Context, d db.Database, artistID, songID string, update map[string]any) (any, error) {
	return d.Update(ctx, SongsTable, map[string]any{"artistid": artistID, "songid": songID}, map[string]any{"$set": update})
}

func SQLDeleteArtistSong(ctx context.Context, d db.Database, artistID, songID string) error {
	_, err := d.Delete(ctx, SongsTable, map[string]any{"artistid": artistID, "songid": songID})
	return err
}
