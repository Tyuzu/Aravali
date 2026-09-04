package musicon

import (
	"context"
	"net/http"
	"scav/infra"
	"scav/infra/db"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
		limit, page = sanitizePagination(limit, page)

		opts := db.FindManyOptions{
			Limit: limit,
			Skip:  (page - 1) * limit,
			Sort:  []bson.E{{Key: "plays", Value: -1}, {Key: "_id", Value: -1}},
		}

		filter := map[string]any{"published": true}

		songs := []Song{}
		if err := app.DB.FindManyWithOptions(ctx, songsCollection, filter, opts, &songs); err != nil {
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
		limit, page = sanitizePagination(limit, page)

		opts := db.FindManyOptions{
			Limit: limit,
			Skip:  (page - 1) * limit,
			Sort:  []bson.E{{Key: "release_date", Value: -1}, {Key: "_id", Value: -1}},
		}

		filter := map[string]any{"published": true}

		albums := []Album{}
		if err := app.DB.FindManyWithOptions(ctx, albumsCollection, filter, opts, &albums); err != nil {
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

		basedOn := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("based_on")))

		filter := map[string]any{"published": true}
		sort := []bson.E{{Key: "_id", Value: -1}} // default stable sort

		switch basedOn {

		case "recently_played":
			filter["plays"] = map[string]any{"$gt": 0}
			sort = []bson.E{{Key: "plays", Value: -1}, {Key: "_id", Value: -1}}

		case "language_en":
			filter["language"] = "en"

		case "genre_pop":
			filter["genre"] = "Pop"
		}

		limit, page := getPaginationParams(r)
		limit, page = sanitizePagination(limit, page)

		opts := db.FindManyOptions{
			Limit: limit,
			Skip:  (page - 1) * limit,
			Sort:  sort,
		}

		songs := []Song{}
		if err := app.DB.FindManyWithOptions(ctx, songsCollection, filter, opts, &songs); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to fetch recommendations")
			return
		}

		respondJSON(w, http.StatusOK, songs, "Personalized recommendations fetched")
	}
}
