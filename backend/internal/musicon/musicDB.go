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
	songsCollection     = config.Collections.SongsCollection
	albumsCollection    = config.Collections.AlbumsCollection
	playlistsCollection = config.Collections.PlaylistsCollection
)

func fetchSongsByIDs(ctx context.Context, ids []string, app *infra.Deps) ([]Song, error) {
	if len(ids) == 0 {
		return []Song{}, nil
	}

	filter := map[string]any{
		"songid":    map[string]any{"$in": ids},
		"published": true,
	}

	var songs []Song
	if err := app.DB.FindMany(ctx, songsCollection, filter, &songs); err != nil {
		return nil, err
	}
	return songs, nil
}

func getPublishedAlbums(ctx context.Context, app *infra.Deps) ([]Album, error) {
	var albums []Album
	if err := app.DB.FindMany(ctx, albumsCollection, map[string]any{"published": true}, &albums); err != nil {
		return nil, err
	}
	return albums, nil
}

func findAlbumByID(ctx context.Context, app *infra.Deps, albumID string) (Album, error) {
	var album Album
	if err := app.DB.FindOne(ctx, albumsCollection, map[string]any{"albumid": albumID}, &album); err != nil {
		return Album{}, err
	}
	return album, nil
}

func findPlaylistByID(ctx context.Context, app *infra.Deps, playlistID string) (Playlist, error) {
	var playlist Playlist
	if err := app.DB.FindOne(ctx, playlistsCollection, map[string]any{"playlistid": playlistID}, &playlist); err != nil {
		return Playlist{}, err
	}
	return playlist, nil
}

func findArtistSongs(ctx context.Context, app *infra.Deps, artistID string, limit, page int) ([]Song, error) {
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
	if err := app.DB.FindManyWithOptions(ctx, songsCollection, filter, opts, &songs); err != nil {
		return nil, err
	}
	return songs, nil
}

func findUserPlaylists(ctx context.Context, app *infra.Deps, userID string) ([]Playlist, error) {
	filter := map[string]any{
		"userid": userID,
		"playlistid": map[string]any{
			"$ne": "likes_" + userID,
		},
	}
	var playlists []Playlist
	if err := app.DB.FindMany(ctx, playlistsCollection, filter, &playlists); err != nil {
		return nil, err
	}
	return playlists, nil
}

func insertPlaylist(ctx context.Context, app *infra.Deps, playlist Playlist) error {
	return app.DB.Insert(ctx, playlistsCollection, playlist)
}

func deletePlaylistForUser(ctx context.Context, app *infra.Deps, playlistID, userID string) (bool, error) {
	filter := map[string]any{
		"playlistid": playlistID,
		"userid":     userID,
	}
	_, err := app.DB.DeleteOne(ctx, playlistsCollection, filter)
	if err != nil {
		return false, err
	}
	return true, nil
}

func addSongToPlaylist(ctx context.Context, app *infra.Deps, playlistID, userID, songID string) error {
	filter := map[string]any{
		"playlistid": playlistID,
		"userid":     userID,
	}
	update := map[string]any{
		"$addToSet": map[string]any{"songs": songID},
		"$set":      map[string]any{"updatedAt": time.Now()},
	}
	_, err := app.DB.UpdateOne(ctx, playlistsCollection, filter, update)
	return err
}

func removeSongFromPlaylist(ctx context.Context, app *infra.Deps, playlistID, userID, songID string) error {
	filter := map[string]any{
		"playlistid": playlistID,
		"userid":     userID,
	}
	update := map[string]any{
		"$pull": map[string]any{"songs": songID},
		"$set":  map[string]any{"updatedAt": time.Now()},
	}
	_, err := app.DB.UpdateOne(ctx, playlistsCollection, filter, update)
	return err
}

func updatePlaylistMeta(ctx context.Context, app *infra.Deps, playlistID, userID, name, description, coverURL string) error {
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
	_, err := app.DB.UpdateOne(ctx, playlistsCollection, filter, update)
	return err
}

func upsertLikedSongsPlaylist(ctx context.Context, app *infra.Deps, userID, songID string) error {
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
	return app.DB.Upsert(ctx, playlistsCollection, filter, update)
}

func unlikeSongFromLikes(ctx context.Context, app *infra.Deps, userID, songID string) error {
	filter := map[string]any{
		"playlistid": "likes_" + userID,
		"userid":     userID,
	}
	update := map[string]any{
		"$pull": map[string]any{"songs": songID},
		"$set":  map[string]any{"updatedAt": time.Now()},
	}
	_, err := app.DB.UpdateOne(ctx, playlistsCollection, filter, update)
	return err
}

func getUserLikedSongs(ctx context.Context, app *infra.Deps, userID string) (Playlist, error) {
	playlistID := fmt.Sprintf("likes_%s", userID)
	var playlist Playlist
	err := app.DB.FindOne(ctx, playlistsCollection, map[string]any{
		"playlistid": playlistID,
		"userid":     userID,
	}, &playlist)
	return playlist, err
}

func getRecommendedSongsList(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions) ([]Song, error) {
	var songs []Song
	if err := app.DB.FindManyWithOptions(ctx, songsCollection, filter, opts, &songs); err != nil {
		return nil, err
	}
	return songs, nil
}

func getRecommendedAlbumsList(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions) ([]Album, error) {
	var albums []Album
	if err := app.DB.FindManyWithOptions(ctx, albumsCollection, filter, opts, &albums); err != nil {
		return nil, err
	}
	return albums, nil
}
