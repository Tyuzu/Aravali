// File: internal/artists/songs/songsSQLDB.go

package songs

import (
	"context"
	"time"

	"scav/config"
	"scav/infra"
	"scav/internal/artists"
)

var (
	SongsTable = config.Tables.SongsTable
)

func ListPublishedSongsByArtist(ctx context.Context, app *infra.Deps, artistID string, result *[]ArtistSong) error {
	query := "artistid = $1 AND published = true"
	args := []any{artistID}
	return app.SQLDB.FindMany(ctx, SongsTable, query, args, result)
}

func FindSongsByArtist(ctx context.Context, app *infra.Deps, artistID string, result *[]ArtistSong) error {
	return ListPublishedSongsByArtist(ctx, app, artistID, result)
}

func FindSongByArtistAndID(ctx context.Context, app *infra.Deps, artistID, songID string, result *ArtistSong) error {
	query := "artistid = $1 AND songid = $2"
	args := []any{artistID, songID}
	return app.SQLDB.FindOne(ctx, SongsTable, query, args, result)
}

func SaveSong(ctx context.Context, app *infra.Deps, song *ArtistSong) error {
	return app.SQLDB.Insert(ctx, SongsTable, song)
}

func InsertArtistSong(ctx context.Context, app *infra.Deps, song *ArtistSong) error {
	return SaveSong(ctx, app, song)
}

func UpdateArtistSongFromPayload(ctx context.Context, app *infra.Deps, artistID, songID string, payload songPayload) (int64, error) {
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
		return 0, artists.ErrNoFieldsToUpdate
	}

	updateFields["updatedAt"] = time.Now()
	return UpdateArtistSong(ctx, app, artistID, songID, updateFields)
}

func UpdateArtistSong(ctx context.Context, app *infra.Deps, artistID, songID string, update map[string]any) (int64, error) {
	query := "artistid = $1 AND songid = $2"
	args := []any{artistID, songID}
	return app.SQLDB.UpdateOne(ctx, SongsTable, query, args, update)
}

func DeleteArtistSong(ctx context.Context, app *infra.Deps, artistID, songID string) error {
	query := "artistid = $1 AND songid = $2"
	_, err := app.SQLDB.DeleteOne(ctx, SongsTable, query, []any{artistID, songID})
	return err
}
