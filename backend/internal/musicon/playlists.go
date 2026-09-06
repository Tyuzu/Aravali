package musicon

import (
	"context"
	"encoding/json"
	"net/http"
	"scav/infra"
	"scav/utils"
	"time"
)

// --------------------------- Playlist Handlers ---------------------------

func GetUserPlaylists(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized or missing user ID")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		playlists, err := findUserPlaylists(ctx, app, userID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to fetch playlists")
			return
		}

		respondJSON(w, http.StatusOK, playlists, "Playlists fetched successfully")
	}
}

func CreatePlaylist(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized or missing user ID")
			return
		}

		type Req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}

		var req Req
		if err := utils.ParseJSON(r, &req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid JSON input")
			return
		}

		if len(req.Name) == 0 || len(req.Name) > 100 {
			respondError(w, http.StatusBadRequest, "Playlist name must be 1-100 characters")
			return
		}

		now := time.Now()

		newPlaylist := Playlist{
			PlaylistID:    "pl_" + utils.GenerateRandomString(12),
			UserID:        userID,
			Name:          req.Name,
			Description:   req.Description,
			Songs:         []string{},
			Duration:      0,
			IsCompilation: false,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := insertPlaylist(ctx, app, newPlaylist); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to create playlist")
			return
		}

		respondJSON(w, http.StatusCreated, newPlaylist, "Playlist created successfully")
	}
}

func DeletePlaylist(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized or missing user ID")
			return
		}

		playlistID := utils.GetParam(r, "playlistid")

		// Prevent deletion of special likes playlist
		if playlistID == "likes_"+userID {
			respondError(w, http.StatusForbidden, "Cannot delete liked songs playlist")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if _, err := deletePlaylistForUser(ctx, app, playlistID, userID); err != nil {
			respondError(w, http.StatusNotFound, "Playlist not found or unauthorized")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"playlist_id": playlistID,
		}, "Playlist deleted successfully")
	}
}

func AddSongToPlaylist(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized or missing user ID")
			return
		}

		playlistID := utils.GetParam(r, "playlistid")

		// Prevent manual modification of likes playlist
		if playlistID == "likes_"+userID {
			respondError(w, http.StatusForbidden, "Liked songs playlist cannot be modified directly")
			return
		}

		var body struct {
			SongID string `json:"songid"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		if body.SongID == "" {
			respondError(w, http.StatusBadRequest, "Missing song ID")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := addSongToPlaylist(ctx, app, playlistID, userID, body.SongID); err != nil {
			respondError(w, http.StatusForbidden, "Playlist not found or unauthorized")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"playlist_id": playlistID,
			"song_id":     body.SongID,
		}, "Song added to playlist")
	}
}

func RemoveSongFromPlaylist(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized or missing user ID")
			return
		}

		playlistID := utils.GetParam(r, "playlistid")
		songID := utils.GetParam(r, "songid")

		// Prevent manual modification of likes playlist
		if playlistID == "likes_"+userID {
			respondError(w, http.StatusForbidden, "Liked songs playlist cannot be modified directly")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := removeSongFromPlaylist(ctx, app, playlistID, userID, songID); err != nil {
			respondError(w, http.StatusForbidden, "Playlist not found or unauthorized")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"playlist_id": playlistID,
			"song_id":     songID,
		}, "Song removed from playlist")
	}
}

func UpdatePlaylistInfo(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized or missing user ID")
			return
		}

		playlistID := utils.GetParam(r, "playlistid")

		// Prevent editing of likes playlist
		if playlistID == "likes_"+userID {
			respondError(w, http.StatusForbidden, "Liked songs playlist cannot be modified")
			return
		}

		type Req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			CoverURL    string `json:"coverUrl"`
		}

		var req Req
		if err := utils.ParseJSON(r, &req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid JSON input")
			return
		}

		if len(req.Name) == 0 || len(req.Name) > 100 {
			respondError(w, http.StatusBadRequest, "Playlist name must be 1-100 characters")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := updatePlaylistMeta(ctx, app, playlistID, userID, req.Name, req.Description, req.CoverURL); err != nil {
			respondError(w, http.StatusForbidden, "Playlist not found or unauthorized")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"playlist_id": playlistID,
		}, "Playlist updated successfully")
	}
}
