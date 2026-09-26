package autocomplete

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/auth"
	"scav/internal/places"
)

var (
	AutocompleteCollection = config.Collections.AutocompleteCollection
)

func findPlacesByQuery(ctx context.Context, app *infra.Deps, query string, places *[]places.Place) error {
	filter := map[string]any{
		"name": map[string]any{
			"$regex":   "^" + query,
			"$options": "i",
		},
	}
	return app.DB.FindMany(ctx, AutocompleteCollection, filter, places)
}

func findUsersByQuery(ctx context.Context, app *infra.Deps, query string, users *[]auth.User) error {
	filter := map[string]any{
		"username": map[string]any{
			"$regex":   "^" + query,
			"$options": "i",
		},
	}
	return app.DB.FindMany(ctx, AutocompleteCollection, filter, users)
}
