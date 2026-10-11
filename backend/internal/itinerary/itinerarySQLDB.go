// File: internal/itinerary/itinerarySQLDB.go

package itinerary

import (
	"context"

	"scav/config"
	"scav/infra"
)

var ItineraryTable = config.Tables.ItineraryTable

func insertItinerary(ctx context.Context, app *infra.Deps, itinerary Itinerary) error {

}

func findItineraryByID(ctx context.Context, app *infra.Deps, itineraryID string) (Itinerary, error) {
	var itinerary Itinerary

	return itinerary, nil
}

func findItineraries(ctx context.Context, app *infra.Deps, query string, args []any) ([]Itinerary, error) {
	var itineraries []Itinerary

	return itineraries, nil
}

func updateItineraryFields(ctx context.Context, app *infra.Deps, itineraryID string, update map[string]any) (int64, error) {

}

func softDeleteItinerary(ctx context.Context, app *infra.Deps, itineraryID, userID string) (int64, error) {

}

func publishItinerary(ctx context.Context, app *infra.Deps, itineraryID, userID string) (int64, error) {

}
