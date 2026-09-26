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
	where := "userid = $1"
	args := []any{userID}

	err := app.SQLDB.FindOne(ctx, followingsTable, where, args, &followData)
	if err != nil {
		return follows.UserFollow{}, err
	}
	return followData, nil
}

func SQLfindSuggestedUsers(ctx context.Context, app *infra.Deps, where string, args []any) ([]UserSuggest, error) {
	var users []UserSuggest
	err := app.SQLDB.FindMany(ctx, usersTable, where, args, &users)
	return users, err
}

func SQLfindNearbyPlaces(ctx context.Context, app *infra.Deps, where string, args []any) ([]places.Place, error) {
	var nearbyPlaces []places.Place
	err := app.SQLDB.FindMany(ctx, placesTable, where, args, &nearbyPlaces)
	return nearbyPlaces, err
}
