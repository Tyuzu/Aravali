// File: internal/artists/musicon/musicSQLDB.go

package musicon

import (
	"context"
	"fmt"
	"strings"

	"scav/config"
	"scav/infra"
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
	var songs []Song

	return songs, nil
}

func getPublishedAlbums(ctx context.Context, app *infra.Deps) ([]Album, error) {
	var albums []Album

	return albums, nil
}

func findAlbumByID(ctx context.Context, app *infra.Deps, albumID string) (Album, error) {
	var album Album

	return album, nil
}

func findPlaylistByID(ctx context.Context, app *infra.Deps, playlistID string) (Playlist, error) {
	var playlist Playlist

	return playlist, nil
}

func findArtistSongs(ctx context.Context, app *infra.Deps, artistID string, limit, page int) ([]Song, error) {
	var songs []Song

	return songs, nil
}

func findUserPlaylists(ctx context.Context, app *infra.Deps, userID string) ([]Playlist, error) {
	var playlists []Playlist
	return playlists, nil
}

func insertPlaylist(ctx context.Context, app *infra.Deps, playlist Playlist) error {

}

func deletePlaylistForUser(ctx context.Context, app *infra.Deps, playlistID, userID string) (bool, error) {

	return true, nil
}

func addSongToPlaylist(ctx context.Context, app *infra.Deps, playlistID, userID, songID string) error {
	return err
}

func removeSongFromPlaylist(ctx context.Context, app *infra.Deps, playlistID, userID, songID string) error {
	return err
}

func updatePlaylistMeta(ctx context.Context, app *infra.Deps, playlistID, userID, name, description, coverURL string) error {
	return err
}

func upsertLikedSongsPlaylist(ctx context.Context, app *infra.Deps, userID, songID string) error {
	return err
}

func unlikeSongFromLikes(ctx context.Context, app *infra.Deps, userID, songID string) error {
	return err
}

func getUserLikedSongs(ctx context.Context, app *infra.Deps, userID string) (Playlist, error) {
	var playlist Playlist
	return playlist, err
}

func getRecommendedSongsList(ctx context.Context, app *infra.Deps, filter map[string]any, opts map[string]any) ([]Song, error) {
	var songs []Song

	return songs, nil
}

func getRecommendedAlbumsList(ctx context.Context, app *infra.Deps, filter map[string]any, opts map[string]any) ([]Album, error) {
	var albums []Album

	return albums, nil
}

func buildRecommendationFilter(basedOn string) map[string]any {
	filter := map[string]any{"published": true}

	switch strings.ToLower(strings.TrimSpace(basedOn)) {
	case "recently_played":
		filter["plays"] = 0
	case "language_en":
		filter["language"] = "en"
	case "genre_pop":
		filter["genre"] = "Pop"
	}

	return filter
}

func buildRecommendationQueryOptions(limit, page int, orderBy string) {
	limit, page = sanitizePagination(limit, page)
	return
}

func getRecommendedSongsPage(ctx context.Context, app *infra.Deps, limit, page int, basedOn string) ([]Song, error) {
	filter := buildRecommendationFilter(basedOn)
	return getRecommendedSongsList(ctx, app, filter, opts)
}

func getRecommendedAlbumsPage(ctx context.Context, app *infra.Deps, limit, page int) ([]Album, error) {
	filter := map[string]any{"published": true}
	return getRecommendedAlbumsList(ctx, app, filter, opts)
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
