// File: internal/artists/musicon/recomm.go

package musicon

import (
	"context"
	"net/http"
	"time"

	"scav/infra"
)

// --------------------------- Helpers ---------------------------

func sanitizePagination(limit, page int) (int, int) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}
	return limit, page
}

// --------------------------- Recommendations ---------------------------

func GetRecommendedSongs(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		limit, page := getPaginationParams(r)
		songs, err := getRecommendedSongsPage(ctx, app, limit, page, "")
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to fetch recommended songs")
			return
		}

		respondJSON(w, http.StatusOK, songs, "Recommended songs fetched")
	}
}

func GetRecommendedAlbums(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		limit, page := getPaginationParams(r)
		albums, err := getRecommendedAlbumsPage(ctx, app, limit, page)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to fetch recommended albums")
			return
		}

		respondJSON(w, http.StatusOK, albums, "Recommended albums fetched")
	}
}

func GetRecommendations(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		basedOn := r.URL.Query().Get("based_on")
		limit, page := getPaginationParams(r)
		songs, err := getRecommendedSongsPage(ctx, app, limit, page, basedOn)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to fetch recommendations")
			return
		}

		respondJSON(w, http.StatusOK, songs, "Personalized recommendations fetched")
	}
}
