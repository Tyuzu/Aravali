package placedb

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var placesTable = config.Tables.PlacesTable
var eventsTable = config.Tables.EventsTable
var productsTable = config.Tables.ProductTable

// Places
func SQLFindPlaces(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindMany(ctx, placesTable, query, args, out)
}

func SQLFindOnePlace(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindOne(ctx, placesTable, query, args, out)
}

func SQLUpdatePlace(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	return app.SQLDB.Update(ctx, placesTable, query, args, update)
}

func SQLInsertPlace(ctx context.Context, app *infra.Deps, place any) error {
	return app.SQLDB.Insert(ctx, placesTable, place)
}

func SQLDeletePlace(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.DeleteOne(ctx, placesTable, query, args)
}

// Events
func SQLCountEvents(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.Count(ctx, eventsTable, query, args)
}

func SQLFindEventsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions, out any) error {
	return app.SQLDB.FindManyWithOptions(ctx, eventsTable, query, args, opts, out)
}

// Place products (generic)
func SQLFindPlaceProducts(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindMany(ctx, productsTable, query, args, out)
}

func SQLInsertPlaceProduct(ctx context.Context, app *infra.Deps, product any) error {
	return app.SQLDB.InsertOne(ctx, productsTable, product)
}

func SQLUpdatePlaceProduct(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	return app.SQLDB.UpdateOne(ctx, productsTable, query, args, update)
}

func SQLDeletePlaceProduct(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.DeleteOne(ctx, productsTable, query, args)
}
