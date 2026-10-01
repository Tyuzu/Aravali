// File: internal/itinerary/GETs_itinerary.go

package itinerary

import (
	"context"
	"net/http"
	"scav/infra"
	"scav/utils"
	"time"
)

// GET /api/itineraries/all/:id
func GetItinerary(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itineraryID := utils.GetParam(r, "id")
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		itinerary, err := findItineraryByID(ctx, app, itineraryID)
		if err != nil {
			http.Error(w, "Itinerary not found", http.StatusNotFound)
			return
		}

		normalizeItinerary(&itinerary)
		utils.RespondWithJSON(w, http.StatusOK, itinerary)
	}
}

// GET /api/itineraries
func GetItineraries(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		itineraries, err := findItineraries(ctx, app, "deleted IS NOT TRUE", nil)
		if err != nil {
			http.Error(w, "Error fetching itineraries", http.StatusInternalServerError)
			return
		}

		if itineraries == nil {
			itineraries = []Itinerary{}
		}

		for i := range itineraries {
			normalizeItinerary(&itineraries[i])
		}

		utils.RespondWithJSON(w, http.StatusOK, itineraries)
	}
}

// GET /api/itineraries/search
func SearchItineraries(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		where := "deleted IS NOT TRUE"
		args := []any{}
		if start := query.Get("start_date"); start != "" {
			where += " AND start_date = $1"
			args = append(args, start)
		}
		if status := query.Get("status"); status != "" {
			where += " AND status = $2"
			args = append(args, status)
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		itineraries, err := findItineraries(ctx, app, where, args)
		if err != nil {
			http.Error(w, "Error fetching itineraries", http.StatusInternalServerError)
			return
		}

		if itineraries == nil {
			itineraries = []Itinerary{}
		}

		for i := range itineraries {
			normalizeItinerary(&itineraries[i])
		}

		utils.RespondWithJSON(w, http.StatusOK, itineraries)
	}
}
