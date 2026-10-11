// File: internal/beats/autocomplete/autocompleteSQLDB.go

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

func findPlacesByQuery(ctx context.Context, app *infra.Deps, query string, places *[]places.Place) error {
	return nil
}

func findUsersByQuery(ctx context.Context, app *infra.Deps, query string, users *[]auth.User) error {
	return nil
}
