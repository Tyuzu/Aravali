package autocomplete

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/auth"
	"scav/internal/places"
)

var (
	AutocompleteTable = config.Tables.AutocompleteTable
)

func SQLfindPlacesByQuery(ctx context.Context, app *infra.Deps, query string, places *[]places.Place) error {
	where := "name ILIKE $1"
	args := []any{query + "%"}

	return app.SQLDB.FindMany(ctx, AutocompleteTable, where, args, places)
}

func SQLfindUsersByQuery(ctx context.Context, app *infra.Deps, query string, users *[]auth.User) error {
	where := "username ILIKE $1"
	args := []any{query + "%"}

	return app.SQLDB.FindMany(ctx, AutocompleteTable, where, args, users)
}
