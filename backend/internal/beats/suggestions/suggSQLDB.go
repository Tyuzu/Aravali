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

func findFollowDataByUserID(_ context.Context, _ *infra.Deps, _ string) (follows.UserFollow, error) {
	var followData follows.UserFollow

	return followData, nil
}

func findSuggestedUsers(_ context.Context, _ *infra.Deps, _ string, _ []any) ([]UserSuggest, error) {
	return nil, nil
}

func findNearbyPlaces(_ context.Context, _ *infra.Deps, _ string, _ []any) ([]places.Place, error) {
	return nil, nil
}
