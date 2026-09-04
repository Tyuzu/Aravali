package autocomplete

import (
	"context"

	"scav/config"
	db "scav/infra/db"
	"scav/internal/auth"
	"scav/internal/places"
)

var (
	AutocompleteCollection = config.Collections.AutocompleteCollection
)

func findPlacesByQuery(ctx context.Context, database db.Database, query string, places *[]places.Place) error {
	filter := map[string]any{
		"name": map[string]any{
			"$regex":   "^" + query,
			"$options": "i",
		},
	}
	return database.FindMany(ctx, AutocompleteCollection, filter, places)
}

func findUsersByQuery(ctx context.Context, database db.Database, query string, users *[]auth.User) error {
	filter := map[string]any{
		"username": map[string]any{
			"$regex":   "^" + query,
			"$options": "i",
		},
	}
	return database.FindMany(ctx, AutocompleteCollection, filter, users)
}
