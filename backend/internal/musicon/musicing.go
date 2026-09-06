package musicon

import (
	"context"
	"fmt"
	"net/http"
	"scav/infra"
	"scav/utils"
	log "scav/utils/logger"
	"strconv"
	"time"
)

// --------------------------- Helpers ---------------------------

func respondJSON(w http.ResponseWriter, status int, data interface{}, message string) {
	utils.RespondWithJSON(w, status, map[string]interface{}{
		"success": true,
		"data":    data,
		"message": message,
	})
}

func respondError(w http.ResponseWriter, status int, message string) {
	utils.RespondWithJSON(w, status, map[string]interface{}{
		"success": false,
		"data":    nil,
		"message": message,
	})
}

func getPaginationParams(r *http.Request) (limit int, page int) {
	limit = 20
	page = 1

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	return
}

// --------------------------- Albums & Songs ---------------------------

func GetAlbums(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		albums, err := getPublishedAlbums(ctx, app)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to fetch albums")
			return
		}

		respondJSON(w, http.StatusOK, albums, "Albums fetched successfully")
	}
}

func GetAlbumSongs(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		albumID := utils.GetParam(r, "albumid")

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		album, err := findAlbumByID(ctx, app, albumID)
		if err != nil {
			respondJSON(w, http.StatusOK, []Song{}, "No songs found for album")
			return
		}

		songs, err := fetchSongsByIDs(ctx, album.Songs, app)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to fetch songs")
			return
		}

		respondJSON(w, http.StatusOK, songs, fmt.Sprintf("Songs for album %s fetched", albumID))
	}
}

func GetPlaylistSongs(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		playlistID := utils.GetParam(r, "playlistid")

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		playlist, err := findPlaylistByID(ctx, app, playlistID)
		if err != nil {
			respondJSON(w, http.StatusOK, []Song{}, "Playlist not found")
			return
		}

		songs, err := fetchSongsByIDs(ctx, playlist.Songs, app)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to fetch songs")
			return
		}

		respondJSON(w, http.StatusOK, songs, fmt.Sprintf("Songs for playlist %s fetched", playlistID))
	}
}

// --------------------------- Artist Songs ---------------------------
func GetArtistsSongs(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		artistID := utils.GetParam(r, "artistid")
		if artistID == "" {
			respondError(w, http.StatusBadRequest, "Missing artist ID")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		limit, page := getPaginationParams(r)
		songs, err := findArtistSongs(ctx, app, artistID, limit, page)
		if err != nil {
			log.Printf("GetArtistsSongs error: %v", err)
			respondError(w, http.StatusInternalServerError, "Failed to fetch artist songs")
			return
		}

		respondJSON(w, http.StatusOK, songs, fmt.Sprintf("Songs for artist %s fetched", artistID))
	}
}
