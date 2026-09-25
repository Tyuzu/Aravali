package placedb

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/db"
)

var placesTable = config.Tables.PlacesTable
var eventsTable = config.Tables.EventsTable
var productsTable = config.Tables.ProductTable

// Places
func SQLFindPlaces(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindMany(ctx, placesTable, filter, out)
}

func SQLFindOnePlace(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindOne(ctx, placesTable, filter, out)
}

func SQLUpdatePlace(ctx context.Context, app *infra.Deps, filter map[string]any, update any) (any, error) {
	return app.DB.Update(ctx, placesTable, filter, update)
}

func SQLInsertPlace(ctx context.Context, app *infra.Deps, place any) error {
	return app.DB.Insert(ctx, placesTable, place)
}

func SQLDeletePlace(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	return app.DB.DeleteOne(ctx, placesTable, filter)
}

// Events
func SQLCountEvents(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	return app.DB.Count(ctx, eventsTable, filter)
}

func SQLFindEventsWithOptions(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions, out any) error {
	return app.DB.FindManyWithOptions(ctx, eventsTable, filter, opts, out)
}

// Place products (generic)
func SQLFindPlaceProducts(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindMany(ctx, productsTable, filter, out)
}

func SQLInsertPlaceProduct(ctx context.Context, app *infra.Deps, product any) error {
	return app.DB.InsertOne(ctx, productsTable, product)
}

func SQLUpdatePlaceProduct(ctx context.Context, app *infra.Deps, filter map[string]any, update any) (any, error) {
	return app.DB.UpdateOne(ctx, productsTable, filter, update)
}

func SQLDeletePlaceProduct(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	return app.DB.DeleteOne(ctx, productsTable, filter)
}
