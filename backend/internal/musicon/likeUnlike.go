package musicon

import (
	"context"
	"net/http"
	"scav/infra"
	"scav/utils"
	log "scav/utils/logger"
	"time"
)

func LikeSong(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		songID := utils.GetParam(r, "songid")
		if songID == "" {
			respondError(w, http.StatusBadRequest, "Missing song ID")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		err := upsertLikedSongsPlaylist(ctx, app, userID, songID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to like song")
			return
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"song_id": songID,
			"liked":   true,
		}, "Song liked successfully")
	}
}
func UnlikeSong(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		songID := utils.GetParam(r, "songid")
		if songID == "" {
			respondError(w, http.StatusBadRequest, "Missing song ID")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		err := unlikeSongFromLikes(ctx, app, userID, songID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to unlike song")
			return
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"song_id": songID,
			"liked":   false,
		}, "Song unliked successfully")
	}
}

// --------------------------- User Likes ---------------------------

func GetUserLikes(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			respondError(w, http.StatusUnauthorized, "Unauthorized or missing user ID")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		playlist, err := getUserLikedSongs(ctx, app, userID)
		if err != nil || len(playlist.Songs) == 0 {
			respondJSON(w, http.StatusOK, []Song{}, "No liked songs found")
			return
		}

		songs, err := fetchSongsByIDs(ctx, playlist.Songs, app)
		if err != nil {
			log.Printf("GetUserLikes fetch error: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to fetch liked songs")
			return
		}

		respondJSON(w, http.StatusOK, songs, "Liked songs fetched successfully")
	}
}
