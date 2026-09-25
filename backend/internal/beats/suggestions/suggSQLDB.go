package suggestions

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/beats/follows"
	"scav/internal/places"
)

var followingsTable = config.Tables.FollowingsTable
var usersTable = config.Tables.UserTable
var placesTable = config.Tables.PlacesTable

func SQLfindFollowDataByUserID(ctx context.Context, app *infra.Deps, userID string) (follows.UserFollow, error) {
	var followData follows.UserFollow
	err := app.DB.FindOne(ctx, followingsTable, map[string]any{"userid": userID}, &followData)
	if err != nil {
		return follows.UserFollow{}, err
	}
	return followData, nil
}

func SQLfindSuggestedUsers(ctx context.Context, app *infra.Deps, filter map[string]any) ([]UserSuggest, error) {
	var users []UserSuggest
	err := app.DB.FindMany(ctx, usersTable, filter, &users)
	return users, err
}

func SQLfindNearbyPlaces(ctx context.Context, app *infra.Deps, filter map[string]any) ([]places.Place, error) {
	var nearbyPlaces []places.Place
	err := app.DB.FindMany(ctx, placesTable, filter, &nearbyPlaces)
	return nearbyPlaces, err
}
