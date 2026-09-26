package itinerary

import (
	"context"

	"scav/config"
	"scav/infra"
)

var ItineraryTable = config.Tables.ItineraryTable

func SQLinsertItinerary(ctx context.Context, app *infra.Deps, itinerary Itinerary) error {
	return app.SQLDB.InsertOne(ctx, ItineraryTable, itinerary)
}

func SQLfindItineraryByID(ctx context.Context, app *infra.Deps, itineraryID string) (Itinerary, error) {
	var itinerary Itinerary
	query := "itineraryid = $1 AND deleted IS NOT TRUE"
	args := []any{itineraryID}

	err := app.SQLDB.FindOne(ctx, ItineraryTable, query, args, &itinerary)
	if err != nil {
		return Itinerary{}, err
	}
	return itinerary, nil
}

func SQLfindItineraries(ctx context.Context, app *infra.Deps, query string, args []any) ([]Itinerary, error) {
	var itineraries []Itinerary
	err := app.SQLDB.FindMany(ctx, ItineraryTable, query, args, &itineraries)
	if err != nil {
		return nil, err
	}
	return itineraries, nil
}

func SQLupdateItineraryFields(ctx context.Context, app *infra.Deps, itineraryID string, update map[string]any) (int64, error) {
	query := "itineraryid = $1"
	args := []any{itineraryID}

	return app.SQLDB.UpdateOne(ctx, ItineraryTable, query, args, update)
}

func SQLsoftDeleteItinerary(ctx context.Context, app *infra.Deps, itineraryID, userID string) (int64, error) {
	query := "itineraryid = $1 AND userid = $2"
	args := []any{itineraryID, userID}
	update := map[string]any{"deleted": true}

	return app.SQLDB.UpdateOne(ctx, ItineraryTable, query, args, update)
}

func SQLpublishItinerary(ctx context.Context, app *infra.Deps, itineraryID, userID string) (int64, error) {
	query := "itineraryid = $1 AND userid = $2"
	args := []any{itineraryID, userID}
	update := map[string]any{"published": true}

	return app.SQLDB.UpdateOne(ctx, ItineraryTable, query, args, update)
}
