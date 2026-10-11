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
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	if len(ids) == 0 {
		return []Song{}, nil
	}
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE songid = ANY($1)", songsTable), ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return []Song{}, nil
}

func getPublishedAlbums(ctx context.Context, app *infra.Deps) ([]Album, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE published = true ORDER BY createdat DESC", albumsTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return []Album{}, nil
}

func findAlbumByID(ctx context.Context, app *infra.Deps, albumID string) (Album, error) {
	if app == nil || app.SQLDB == nil {
		return Album{}, nil
	}
	if albumID == "" {
		return Album{}, nil
	}
	row := app.SQLDB.QueryRow(ctx, fmt.Sprintf("SELECT * FROM %s WHERE albumid = $1 LIMIT 1", albumsTable), albumID)
	if row == nil {
		return Album{}, nil
	}
	return Album{}, nil
}

func findPlaylistByID(ctx context.Context, app *infra.Deps, playlistID string) (Playlist, error) {
	if app == nil || app.SQLDB == nil {
		return Playlist{}, nil
	}
	if playlistID == "" {
		return Playlist{}, nil
	}
	row := app.SQLDB.QueryRow(ctx, fmt.Sprintf("SELECT * FROM %s WHERE playlistid = $1 LIMIT 1", playlistsTable), playlistID)
	if row == nil {
		return Playlist{}, nil
	}
	return Playlist{}, nil
}

func findArtistSongs(ctx context.Context, app *infra.Deps, artistID string, limit, page int) ([]Song, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	if artistID == "" {
		return []Song{}, nil
	}
	limit, page = sanitizePagination(limit, page)
	offset := (page - 1) * limit
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE artistid = $1 ORDER BY uploadedat DESC LIMIT $2 OFFSET $3", songsTable), artistID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return []Song{}, nil
}

func findUserPlaylists(ctx context.Context, app *infra.Deps, userID string) ([]Playlist, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	if userID == "" {
		return []Playlist{}, nil
	}
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE userid = $1 ORDER BY updatedat DESC", playlistsTable), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return []Playlist{}, nil
}

func insertPlaylist(ctx context.Context, app *infra.Deps, playlist Playlist) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, fmt.Sprintf("INSERT INTO %s (playlistid, userid, name, description, songs, createdat, updatedat, duration, iscompilation, copyrights, coverurl) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", playlistsTable), playlist.PlaylistID, playlist.UserID, playlist.Name, playlist.Description, playlist.Songs, playlist.CreatedAt, playlist.UpdatedAt, playlist.Duration, playlist.IsCompilation, playlist.Copyrights, playlist.CoverURL)
	return err
}

func deletePlaylistForUser(ctx context.Context, app *infra.Deps, playlistID, userID string) (bool, error) {
	if app == nil || app.SQLDB == nil {
		return true, nil
	}
	if playlistID == "" || userID == "" {
		return false, nil
	}
	res, err := app.SQLDB.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE playlistid = $1 AND userid = $2", playlistsTable), playlistID, userID)
	if err != nil {
		return false, err
	}
	return res.RowsAffected() > 0, nil
}

func addSongToPlaylist(ctx context.Context, app *infra.Deps, playlistID, userID, songID string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if playlistID == "" || userID == "" || songID == "" {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, fmt.Sprintf("UPDATE %s SET songs = COALESCE(songs, '{}'::text[]) || $1, updatedat = NOW() WHERE playlistid = $2 AND userid = $3", playlistsTable), []string{songID}, playlistID, userID)
	return err
}

func removeSongFromPlaylist(ctx context.Context, app *infra.Deps, playlistID, userID, songID string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if playlistID == "" || userID == "" || songID == "" {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, fmt.Sprintf("UPDATE %s SET songs = array_remove(songs, $1), updatedat = NOW() WHERE playlistid = $2 AND userid = $3", playlistsTable), songID, playlistID, userID)
	return err
}

func updatePlaylistMeta(ctx context.Context, app *infra.Deps, playlistID, userID, name, description, coverURL string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if playlistID == "" || userID == "" {
		return nil
	}
	setParts := []string{"updatedat = NOW()"}
	args := []any{playlistID, userID}
	idx := len(args) + 1
	if strings.TrimSpace(name) != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", idx))
		args = append(args, name)
		idx++
	}
	if strings.TrimSpace(description) != "" {
		setParts = append(setParts, fmt.Sprintf("description = $%d", idx))
		args = append(args, description)
		idx++
	}
	if strings.TrimSpace(coverURL) != "" {
		setParts = append(setParts, fmt.Sprintf("coverurl = $%d", idx))
		args = append(args, coverURL)
	}
	_, err := app.SQLDB.Exec(ctx, fmt.Sprintf("UPDATE %s SET %s WHERE playlistid = $1 AND userid = $2", playlistsTable, strings.Join(setParts, ", ")), args...)
	return err
}

func upsertLikedSongsPlaylist(ctx context.Context, app *infra.Deps, userID, songID string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if userID == "" || songID == "" {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, fmt.Sprintf("INSERT INTO %s (playlistid, userid, name, songs, createdat, updatedat) VALUES ($1, $2, $3, $4, NOW(), NOW()) ON CONFLICT (userid, name) DO UPDATE SET songs = COALESCE(%s.songs, '{}'::text[]) || $4, updatedat = NOW()", playlistsTable, playlistsTable), userID+"-liked", userID, "liked_songs", []string{songID})
	return err
}

func unlikeSongFromLikes(ctx context.Context, app *infra.Deps, userID, songID string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if userID == "" || songID == "" {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, fmt.Sprintf("UPDATE %s SET songs = array_remove(songs, $1), updatedat = NOW() WHERE userid = $2 AND name = 'liked_songs'", playlistsTable), songID, userID)
	return err
}

func getUserLikedSongs(ctx context.Context, app *infra.Deps, userID string) (Playlist, error) {
	if app == nil || app.SQLDB == nil {
		return Playlist{}, nil
	}
	if userID == "" {
		return Playlist{}, nil
	}
	row := app.SQLDB.QueryRow(ctx, fmt.Sprintf("SELECT * FROM %s WHERE userid = $1 AND name = 'liked_songs' LIMIT 1", playlistsTable), userID)
	if row == nil {
		return Playlist{}, nil
	}
	return Playlist{}, nil
}

func getRecommendedSongsList(ctx context.Context, app *infra.Deps, filter map[string]any, opts map[string]any) ([]Song, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	query, args := buildFilterQuery(filter)
	limit := 20
	page := 1
	if opts != nil {
		if v, ok := opts["limit"].(int); ok {
			limit = v
		}
		if v, ok := opts["page"].(int); ok {
			page = v
		}
	}
	offset := (page - 1) * limit
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE %s ORDER BY uploadedat DESC LIMIT $%d OFFSET $%d", songsTable, query, len(args)+1, len(args)+2), append(args, limit, offset)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return []Song{}, nil
}

func getRecommendedAlbumsList(ctx context.Context, app *infra.Deps, filter map[string]any, opts map[string]any) ([]Album, error) {
	if app == nil || app.SQLDB == nil {
		return nil, nil
	}
	query, args := buildFilterQuery(filter)
	limit := 20
	page := 1
	if opts != nil {
		if v, ok := opts["limit"].(int); ok {
			limit = v
		}
		if v, ok := opts["page"].(int); ok {
			page = v
		}
	}
	offset := (page - 1) * limit
	rows, err := app.SQLDB.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE %s ORDER BY createdat DESC LIMIT $%d OFFSET $%d", albumsTable, query, len(args)+1, len(args)+2), append(args, limit, offset)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return []Album{}, nil
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

func buildRecommendationQueryOptions(limit, page int, orderBy string) map[string]any {
	limit, page = sanitizePagination(limit, page)
	if orderBy == "" {
		orderBy = "created_at DESC"
	}
	return map[string]any{
		"limit":   limit,
		"page":    page,
		"orderBy": orderBy,
	}
}

func getRecommendedSongsPage(ctx context.Context, app *infra.Deps, limit, page int, basedOn string) ([]Song, error) {
	filter := buildRecommendationFilter(basedOn)
	opts := buildRecommendationQueryOptions(limit, page, "plays DESC")
	return getRecommendedSongsList(ctx, app, filter, opts)
}

func getRecommendedAlbumsPage(ctx context.Context, app *infra.Deps, limit, page int) ([]Album, error) {
	filter := map[string]any{"published": true}
	opts := buildRecommendationQueryOptions(limit, page, "created_at DESC")
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
