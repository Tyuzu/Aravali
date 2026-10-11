// File: internal/artists/songs/songsSQLDB.go

package songs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/internal/artists"
)

var (
	SongsTable = config.Tables.SongsTable
)

func ListPublishedSongsByArtist(ctx context.Context, app *infra.Deps, artistID string, result *[]ArtistSong) error {
	if app == nil || app.SQLDB == nil || result == nil {
		return nil
	}
	if artistID == "" {
		*result = nil
		return nil
	}
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE artistid = $1 AND published = true ORDER BY uploadedat DESC", SongsTable), artistID)
	if err != nil {
		return err
	}
	defer rows.Close()
	*result = nil
	return nil
}

func FindSongsByArtist(ctx context.Context, app *infra.Deps, artistID string, result *[]ArtistSong) error {
	return ListPublishedSongsByArtist(ctx, app, artistID, result)
}

func FindSongByArtistAndID(ctx context.Context, app *infra.Deps, artistID, songID string, result *ArtistSong) error {
	if app == nil || app.SQLDB == nil || result == nil {
		return nil
	}
	if artistID == "" || songID == "" {
		*result = ArtistSong{}
		return nil
	}
	row := app.SQLDB.QueryRow(ctx, fmt.Sprintf("SELECT * FROM %s WHERE artistid = $1 AND songid = $2 LIMIT 1", SongsTable), artistID, songID)
	if row == nil {
		*result = ArtistSong{}
		return nil
	}
	*result = ArtistSong{}
	return nil
}

func SaveSong(ctx context.Context, app *infra.Deps, song *ArtistSong) error {
	if app == nil || app.SQLDB == nil || song == nil {
		return nil
	}
	if song.UploadedAt.IsZero() {
		song.UploadedAt = time.Now().UTC()
	}
	_, err := app.SQLDB.Exec(ctx, fmt.Sprintf("INSERT INTO %s (songid, artistid, title, genre, duration, description, audiourl, poster, audioextn, posterextn, published, uploadedat, language) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", SongsTable), song.SongID, song.ArtistID, song.Title, song.Genre, song.Duration, song.Description, song.AudioURL, song.Poster, song.AudioExtn, song.PosterExtn, song.Published, song.UploadedAt, song.Language)
	return err
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
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if artistID == "" || songID == "" || len(update) == 0 {
		return 0, nil
	}
	setParts := make([]string, 0, len(update))
	args := make([]any, 0, len(update))
	idx := 1
	for key, value := range update {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", key, idx))
		args = append(args, value)
		idx++
	}
	res, err := app.SQLDB.Exec(ctx, fmt.Sprintf("UPDATE %s SET %s WHERE artistid = $%d AND songid = $%d", SongsTable, joinSetParts(setParts), len(args)+1, len(args)+2), append(args, artistID, songID)...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func joinSetParts(parts []string) string {
	return strings.Join(parts, ", ")
}

func DeleteArtistSong(ctx context.Context, app *infra.Deps, artistID, songID string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if artistID == "" || songID == "" {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE artistid = $1 AND songid = $2", SongsTable), artistID, songID)
	return err
}
