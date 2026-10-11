// File: internal/itinerary/itinerarySQLDB.go

package itinerary

import (
	"context"

	"scav/config"
	"scav/infra"
)

var ItineraryTable = config.Tables.ItineraryTable

func insertItinerary(ctx context.Context, app *infra.Deps, itinerary Itinerary) error {
	return nil
}

func findItineraryByID(ctx context.Context, app *infra.Deps, itineraryID string) (Itinerary, error) {
	return Itinerary{}, nil
}

func findItineraries(ctx context.Context, app *infra.Deps, query string, args []any) ([]Itinerary, error) {
	return nil, nil
}

func updateItineraryFields(ctx context.Context, app *infra.Deps, itineraryID string, update map[string]any) (int64, error) {
	return 0, nil
}

func softDeleteItinerary(ctx context.Context, app *infra.Deps, itineraryID, userID string) (int64, error) {
	return 0, nil
}

func publishItinerary(ctx context.Context, app *infra.Deps, itineraryID, userID string) (int64, error) {
	return 0, nil
}
