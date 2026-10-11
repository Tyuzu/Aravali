// File: internal/places/placedb/placeSQLDB.go

package placedb

import (
	"context"

	"scav/config"
	"scav/infra"
)

var placesTable = config.Tables.PlacesTable
var eventsTable = config.Tables.EventsTable
var productsTable = config.Tables.ProductTable
var membershipsTable = config.Tables.MembershipsTable

// Places
func FindPlaces(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
}

func FindOnePlace(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
}

func UpdatePlace(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
}

func InsertPlace(ctx context.Context, app *infra.Deps, place any) error {
}

func DeletePlace(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
}

// Events
func CountEvents(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
}

func FindEventsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any, out any) error {
}

// Place products (generic)
func FindPlaceProducts(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
}

func InsertPlaceProduct(ctx context.Context, app *infra.Deps, product any) error {
}

func UpdatePlaceProduct(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
}

func DeletePlaceProduct(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
}

// Memberships
func FindPlaceMemberships(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
}

func FindOneMembership(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
}

func InsertMembership(ctx context.Context, app *infra.Deps, membership any) error {
}

func UpdateMembership(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
}

func DeleteMembership(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
}
