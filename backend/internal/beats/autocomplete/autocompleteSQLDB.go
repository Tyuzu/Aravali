package autocomplete

import (
	"context"

	"scav/config"
	db "scav/infra/db"
	"scav/internal/auth"
	"scav/internal/places"
)

var (
	AutocompleteTable = config.Tables.AutocompleteTable
)

func SQLfindPlacesByQuery(ctx context.Context, database db.Database, query string, places *[]places.Place) error {
	filter := map[string]any{
		"name": map[string]any{
			"$regex":   "^" + query,
			"$options": "i",
		},
	}
	return database.FindMany(ctx, AutocompleteTable, filter, places)
}

func SQLfindUsersByQuery(ctx context.Context, database db.Database, query string, users *[]auth.User) error {
	filter := map[string]any{
		"username": map[string]any{
			"$regex":   "^" + query,
			"$options": "i",
		},
	}
	return database.FindMany(ctx, AutocompleteTable, filter, users)
}
