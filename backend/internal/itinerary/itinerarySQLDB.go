package itinerary

import (
	"context"
	"scav/config"
	"scav/infra"
)

var ItineraryTable = config.Tables.ItineraryTable

func SQLinsertItinerary(ctx context.Context, app *infra.Deps, itinerary Itinerary) error {
	return app.DB.Insert(ctx, ItineraryTable, itinerary)
}

func SQLfindItineraryByID(ctx context.Context, app *infra.Deps, itineraryID string) (Itinerary, error) {
	var itinerary Itinerary
	err := app.DB.FindOne(ctx, ItineraryTable, map[string]any{
		"itineraryid": itineraryID,
		"deleted":     map[string]any{"$ne": true},
	}, &itinerary)
	if err != nil {
		return Itinerary{}, err
	}
	return itinerary, nil
}

func SQLfindItineraries(ctx context.Context, app *infra.Deps, filter map[string]any) ([]Itinerary, error) {
	var itineraries []Itinerary
	err := app.DB.FindMany(ctx, ItineraryTable, filter, &itineraries)
	if err != nil {
		return nil, err
	}
	return itineraries, nil
}

func SQLupdateItineraryFields(ctx context.Context, app *infra.Deps, itineraryID string, update map[string]any) (any, error) {
	return app.DB.UpdateOne(ctx, ItineraryTable, map[string]any{"itineraryid": itineraryID}, update)
}

func SQLsoftDeleteItinerary(ctx context.Context, app *infra.Deps, itineraryID, userID string) (any, error) {
	update := map[string]any{"$set": map[string]any{"deleted": true}}
	return app.DB.UpdateOne(ctx, ItineraryTable, map[string]any{"itineraryid": itineraryID, "userid": userID}, update)
}

func SQLpublishItinerary(ctx context.Context, app *infra.Deps, itineraryID, userID string) (any, error) {
	update := map[string]any{"$set": map[string]any{"published": true}}
	return app.DB.UpdateOne(ctx, ItineraryTable, map[string]any{"itineraryid": itineraryID, "userid": userID}, update)
}
