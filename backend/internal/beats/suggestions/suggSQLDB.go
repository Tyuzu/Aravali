// File: internal/beats/suggestions/suggSQLDB.go

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

func findFollowDataByUserID(ctx context.Context, app *infra.Deps, userID string) (follows.UserFollow, error) {
	var followData follows.UserFollow

	return followData, nil
}

func findSuggestedUsers(ctx context.Context, app *infra.Deps, where string, args []any) ([]UserSuggest, error) {
	return nil, nil
}

func findNearbyPlaces(ctx context.Context, app *infra.Deps, where string, args []any) ([]places.Place, error) {
	return nil, nil
}
