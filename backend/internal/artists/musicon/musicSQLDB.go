package musicon

import (
	"context"
	"fmt"
	"scav/config"
	"scav/infra"
	"scav/infra/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	songsTable     = config.Tables.SongsTable
	albumsTable    = config.Tables.AlbumsTable
	playlistsTable = config.Tables.PlaylistsTable
)

func SQLfetchSongsByIDs(ctx context.Context, ids []string, app *infra.Deps) ([]Song, error) {
	if len(ids) == 0 {
		return []Song{}, nil
	}

	filter := map[string]any{
		"songid":    map[string]any{"$in": ids},
		"published": true,
	}

	var songs []Song
	if err := app.DB.FindMany(ctx, songsTable, filter, &songs); err != nil {
		return nil, err
	}
	return songs, nil
}

func SQLgetPublishedAlbums(ctx context.Context, app *infra.Deps) ([]Album, error) {
	var albums []Album
	if err := app.DB.FindMany(ctx, albumsTable, map[string]any{"published": true}, &albums); err != nil {
		return nil, err
	}
	return albums, nil
}

func SQLfindAlbumByID(ctx context.Context, app *infra.Deps, albumID string) (Album, error) {
	var album Album
	if err := app.DB.FindOne(ctx, albumsTable, map[string]any{"albumid": albumID}, &album); err != nil {
		return Album{}, err
	}
	return album, nil
}

func SQLfindPlaylistByID(ctx context.Context, app *infra.Deps, playlistID string) (Playlist, error) {
	var playlist Playlist
	if err := app.DB.FindOne(ctx, playlistsTable, map[string]any{"playlistid": playlistID}, &playlist); err != nil {
		return Playlist{}, err
	}
	return playlist, nil
}

func SQLfindArtistSongs(ctx context.Context, app *infra.Deps, artistID string, limit, page int) ([]Song, error) {
	skip := (page - 1) * limit
	filter := map[string]any{
		"artistid":  artistID,
		"published": true,
	}
	opts := db.FindManyOptions{
		Limit: limit,
		Skip:  skip,
		Sort:  []bson.E{{Key: "uploadedAt", Value: -1}},
	}
	var songs []Song
	if err := app.DB.FindManyWithOptions(ctx, songsTable, filter, opts, &songs); err != nil {
		return nil, err
	}
	return songs, nil
}

func SQLfindUserPlaylists(ctx context.Context, app *infra.Deps, userID string) ([]Playlist, error) {
	filter := map[string]any{
		"userid": userID,
		"playlistid": map[string]any{
			"$ne": "likes_" + userID,
		},
	}
	var playlists []Playlist
	if err := app.DB.FindMany(ctx, playlistsTable, filter, &playlists); err != nil {
		return nil, err
	}
	return playlists, nil
}

func SQLinsertPlaylist(ctx context.Context, app *infra.Deps, playlist Playlist) error {
	return app.DB.Insert(ctx, playlistsTable, playlist)
}

func SQLdeletePlaylistForUser(ctx context.Context, app *infra.Deps, playlistID, userID string) (bool, error) {
	filter := map[string]any{
		"playlistid": playlistID,
		"userid":     userID,
	}
	_, err := app.DB.DeleteOne(ctx, playlistsTable, filter)
	if err != nil {
		return false, err
	}
	return true, nil
}

func SQLaddSongToPlaylist(ctx context.Context, app *infra.Deps, playlistID, userID, songID string) error {
	filter := map[string]any{
		"playlistid": playlistID,
		"userid":     userID,
	}
	update := map[string]any{
		"$addToSet": map[string]any{"songs": songID},
		"$set":      map[string]any{"updatedAt": time.Now()},
	}
	_, err := app.DB.UpdateOne(ctx, playlistsTable, filter, update)
	return err
}

func SQLremoveSongFromPlaylist(ctx context.Context, app *infra.Deps, playlistID, userID, songID string) error {
	filter := map[string]any{
		"playlistid": playlistID,
		"userid":     userID,
	}
	update := map[string]any{
		"$pull": map[string]any{"songs": songID},
		"$set":  map[string]any{"updatedAt": time.Now()},
	}
	_, err := app.DB.UpdateOne(ctx, playlistsTable, filter, update)
	return err
}

func SQLupdatePlaylistMeta(ctx context.Context, app *infra.Deps, playlistID, userID, name, description, coverURL string) error {
	filter := map[string]any{
		"playlistid": playlistID,
		"userid":     userID,
	}
	update := map[string]any{
		"$set": map[string]any{
			"name":        name,
			"description": description,
			"coverUrl":    coverURL,
			"updatedAt":   time.Now(),
		},
	}
	_, err := app.DB.UpdateOne(ctx, playlistsTable, filter, update)
	return err
}

func SQLupsertLikedSongsPlaylist(ctx context.Context, app *infra.Deps, userID, songID string) error {
	playlistID := "likes_" + userID
	now := time.Now()
	filter := map[string]any{
		"playlistid": playlistID,
		"userid":     userID,
	}
	update := map[string]any{
		"$setOnInsert": map[string]any{
			"playlistid":  playlistID,
			"userid":      userID,
			"name":        "Liked Songs",
			"description": "Auto-generated liked songs playlist",
			"songs":       []string{},
			"duration":    0,
			"createdAt":   now,
		},
		"$addToSet": map[string]any{
			"songs": songID,
		},
		"$set": map[string]any{
			"updatedAt": now,
		},
	}
	return app.DB.Upsert(ctx, playlistsTable, filter, update)
}

func SQLunlikeSongFromLikes(ctx context.Context, app *infra.Deps, userID, songID string) error {
	filter := map[string]any{
		"playlistid": "likes_" + userID,
		"userid":     userID,
	}
	update := map[string]any{
		"$pull": map[string]any{"songs": songID},
		"$set":  map[string]any{"updatedAt": time.Now()},
	}
	_, err := app.DB.UpdateOne(ctx, playlistsTable, filter, update)
	return err
}

func SQLgetUserLikedSongs(ctx context.Context, app *infra.Deps, userID string) (Playlist, error) {
	playlistID := fmt.Sprintf("likes_%s", userID)
	var playlist Playlist
	err := app.DB.FindOne(ctx, playlistsTable, map[string]any{
		"playlistid": playlistID,
		"userid":     userID,
	}, &playlist)
	return playlist, err
}

func SQLgetRecommendedSongsList(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions) ([]Song, error) {
	var songs []Song
	if err := app.DB.FindManyWithOptions(ctx, songsTable, filter, opts, &songs); err != nil {
		return nil, err
	}
	return songs, nil
}

func SQLgetRecommendedAlbumsList(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions) ([]Album, error) {
	var albums []Album
	if err := app.DB.FindManyWithOptions(ctx, albumsTable, filter, opts, &albums); err != nil {
		return nil, err
	}
	return albums, nil
}
