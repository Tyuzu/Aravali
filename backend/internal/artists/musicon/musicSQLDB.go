// File: internal/artists/musicon/musicSQLDB.go

package musicon

import (
	"context"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var (
	songsTable     = config.Tables.SongsTable
	albumsTable    = config.Tables.AlbumsTable
	playlistsTable = config.Tables.PlaylistsTable
)

func fetchSongsByIDs(ctx context.Context, ids []string, app *infra.Deps) ([]Song, error) {
	if len(ids) == 0 {
		return []Song{}, nil
	}
	query := "songid = ANY($1) AND published = true"
	args := []any{ids}
	var songs []Song
	if err := app.SQLDB.FindMany(ctx, songsTable, query, args, &songs); err != nil {
		return nil, err
	}
	return songs, nil
}

func getPublishedAlbums(ctx context.Context, app *infra.Deps) ([]Album, error) {
	var albums []Album
	if err := app.SQLDB.FindMany(ctx, albumsTable, "published = true", nil, &albums); err != nil {
		return nil, err
	}
	return albums, nil
}

func findAlbumByID(ctx context.Context, app *infra.Deps, albumID string) (Album, error) {
	var album Album
	if err := app.SQLDB.FindOne(ctx, albumsTable, "albumid = $1", []any{albumID}, &album); err != nil {
		return Album{}, err
	}
	return album, nil
}

func findPlaylistByID(ctx context.Context, app *infra.Deps, playlistID string) (Playlist, error) {
	var playlist Playlist
	if err := app.SQLDB.FindOne(ctx, playlistsTable, "playlistid = $1", []any{playlistID}, &playlist); err != nil {
		return Playlist{}, err
	}
	return playlist, nil
}

func findArtistSongs(ctx context.Context, app *infra.Deps, artistID string, limit, page int) ([]Song, error) {
	skip := (page - 1) * limit
	query := "artistid = $1 AND published = true"
	args := []any{artistID}
	opts := sqldb.FindManyOptions{Limit: int64(limit), Offset: int64(skip), OrderBy: "uploadedat DESC"}
	var songs []Song
	if err := app.SQLDB.FindManyWithOptions(ctx, songsTable, query, args, opts, &songs); err != nil {
		return nil, err
	}
	return songs, nil
}

func findUserPlaylists(ctx context.Context, app *infra.Deps, userID string) ([]Playlist, error) {
	query := "userid = $1 AND playlistid <> $2"
	args := []any{userID, "likes_" + userID}
	var playlists []Playlist
	if err := app.SQLDB.FindMany(ctx, playlistsTable, query, args, &playlists); err != nil {
		return nil, err
	}
	return playlists, nil
}

func insertPlaylist(ctx context.Context, app *infra.Deps, playlist Playlist) error {
	return app.SQLDB.Insert(ctx, playlistsTable, playlist)
}

func deletePlaylistForUser(ctx context.Context, app *infra.Deps, playlistID, userID string) (bool, error) {
	query := "playlistid = $1 AND userid = $2"
	args := []any{playlistID, userID}
	_, err := app.SQLDB.DeleteOne(ctx, playlistsTable, query, args)
	if err != nil {
		return false, err
	}
	return true, nil
}

func addSongToPlaylist(ctx context.Context, app *infra.Deps, playlistID, userID, songID string) error {
	query := "playlistid = $1 AND userid = $2"
	args := []any{playlistID, userID}
	var playlist Playlist
	if err := app.SQLDB.FindOne(ctx, playlistsTable, query, args, &playlist); err != nil {
		return err
	}
	playlist.Songs = append(playlist.Songs, songID)
	playlist.UpdatedAt = time.Now()
	_, err := app.SQLDB.UpdateOne(ctx, playlistsTable, query, args, map[string]any{"songs": playlist.Songs, "updatedat": playlist.UpdatedAt})
	return err
}

func removeSongFromPlaylist(ctx context.Context, app *infra.Deps, playlistID, userID, songID string) error {
	query := "playlistid = $1 AND userid = $2"
	args := []any{playlistID, userID}
	var playlist Playlist
	if err := app.SQLDB.FindOne(ctx, playlistsTable, query, args, &playlist); err != nil {
		return err
	}
	filtered := make([]string, 0, len(playlist.Songs))
	for _, id := range playlist.Songs {
		if id != songID {
			filtered = append(filtered, id)
		}
	}
	playlist.Songs = filtered
	playlist.UpdatedAt = time.Now()
	_, err := app.SQLDB.UpdateOne(ctx, playlistsTable, query, args, map[string]any{"songs": playlist.Songs, "updatedat": playlist.UpdatedAt})
	return err
}

func updatePlaylistMeta(ctx context.Context, app *infra.Deps, playlistID, userID, name, description, coverURL string) error {
	query := "playlistid = $1 AND userid = $2"
	args := []any{playlistID, userID}
	_, err := app.SQLDB.UpdateOne(ctx, playlistsTable, query, args, map[string]any{
		"name":        name,
		"description": description,
		"coverurl":    coverURL,
		"updatedat":   time.Now(),
	})
	return err
}

func upsertLikedSongsPlaylist(ctx context.Context, app *infra.Deps, userID, songID string) error {
	playlistID := "likes_" + userID
	now := time.Now()
	query := "playlistid = $1 AND userid = $2"
	args := []any{playlistID, userID}
	var playlist Playlist
	if err := app.SQLDB.FindOne(ctx, playlistsTable, query, args, &playlist); err != nil {
		playlist = Playlist{UserID: userID, PlaylistID: playlistID, Name: "Liked Songs", Description: "Auto-generated liked songs playlist", Songs: []string{}, Duration: 0, CreatedAt: now, UpdatedAt: now}
		if err := app.SQLDB.Insert(ctx, playlistsTable, playlist); err != nil {
			return err
		}
	}
	playlist.Songs = append(playlist.Songs, songID)
	playlist.UpdatedAt = now
	_, err := app.SQLDB.UpdateOne(ctx, playlistsTable, query, args, map[string]any{"songs": playlist.Songs, "updatedat": playlist.UpdatedAt})
	return err
}

func unlikeSongFromLikes(ctx context.Context, app *infra.Deps, userID, songID string) error {
	query := "playlistid = $1 AND userid = $2"
	args := []any{"likes_" + userID, userID}
	var playlist Playlist
	if err := app.SQLDB.FindOne(ctx, playlistsTable, query, args, &playlist); err != nil {
		return err
	}
	filtered := make([]string, 0, len(playlist.Songs))
	for _, id := range playlist.Songs {
		if id != songID {
			filtered = append(filtered, id)
		}
	}
	playlist.Songs = filtered
	playlist.UpdatedAt = time.Now()
	_, err := app.SQLDB.UpdateOne(ctx, playlistsTable, query, args, map[string]any{"songs": playlist.Songs, "updatedat": playlist.UpdatedAt})
	return err
}

func getUserLikedSongs(ctx context.Context, app *infra.Deps, userID string) (Playlist, error) {
	playlistID := fmt.Sprintf("likes_%s", userID)
	var playlist Playlist
	err := app.SQLDB.FindOne(ctx, playlistsTable, "playlistid = $1 AND userid = $2", []any{playlistID, userID}, &playlist)
	return playlist, err
}

func getRecommendedSongsList(ctx context.Context, app *infra.Deps, filter map[string]any, opts sqldb.FindManyOptions) ([]Song, error) {
	var songs []Song
	query, args := buildFilterQuery(filter)
	if err := app.SQLDB.FindManyWithOptions(ctx, songsTable, query, args, opts, &songs); err != nil {
		return nil, err
	}
	return songs, nil
}

func getRecommendedAlbumsList(ctx context.Context, app *infra.Deps, filter map[string]any, opts sqldb.FindManyOptions) ([]Album, error) {
	var albums []Album
	query, args := buildFilterQuery(filter)
	if err := app.SQLDB.FindManyWithOptions(ctx, albumsTable, query, args, opts, &albums); err != nil {
		return nil, err
	}
	return albums, nil
}

func buildFilterQuery(filter map[string]any) (string, []any) {
	if len(filter) == 0 {
		return "1 = 1", nil
	}
	clauses := make([]string, 0, len(filter))
	args := make([]any, 0, len(filter))
	for key, value := range filter {
		if value == nil {
			continue
		}
		clauses = append(clauses, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	if len(clauses) == 0 {
		return "1 = 1", nil
	}
	return strings.Join(clauses, " AND "), args
}
