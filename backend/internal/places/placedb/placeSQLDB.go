// File: internal/places/placedb/placeSQLDB.go

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
var membershipsTable = config.Tables.MembershipsTable

// Places
func FindPlaces(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindMany(ctx, placesTable, query, args, out)
}

func FindOnePlace(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindOne(ctx, placesTable, query, args, out)
}

func UpdatePlace(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	return app.SQLDB.Update(ctx, placesTable, query, args, update)
}

func InsertPlace(ctx context.Context, app *infra.Deps, place any) error {
	return app.SQLDB.Insert(ctx, placesTable, place)
}

func DeletePlace(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.DeleteOne(ctx, placesTable, query, args)
}

// Events
func CountEvents(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.Count(ctx, eventsTable, query, args)
}

func FindEventsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions, out any) error {
	return app.SQLDB.FindManyWithOptions(ctx, eventsTable, query, args, opts, out)
}

// Place products (generic)
func FindPlaceProducts(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindMany(ctx, productsTable, query, args, out)
}

func InsertPlaceProduct(ctx context.Context, app *infra.Deps, product any) error {
	return app.SQLDB.InsertOne(ctx, productsTable, product)
}

func UpdatePlaceProduct(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	return app.SQLDB.UpdateOne(ctx, productsTable, query, args, update)
}

func DeletePlaceProduct(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.DeleteOne(ctx, productsTable, query, args)
}

// Memberships
func FindPlaceMemberships(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindMany(ctx, membershipsTable, query, args, out)
}

func FindOneMembership(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindOne(ctx, membershipsTable, query, args, out)
}

func InsertMembership(ctx context.Context, app *infra.Deps, membership any) error {
	return app.SQLDB.InsertOne(ctx, membershipsTable, membership)
}

func UpdateMembership(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	return app.SQLDB.UpdateOne(ctx, membershipsTable, query, args, update)
}

func DeleteMembership(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.DeleteOne(ctx, membershipsTable, query, args)
}
