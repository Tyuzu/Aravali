package suggestions

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/beats/follows"
	"scav/internal/places"
)

var followingsCollection = config.Collections.FollowingsCollection
var usersCollection = config.Collections.UserCollection
var placesCollection = config.Collections.PlacesCollection

func findFollowDataByUserID(ctx context.Context, app *infra.Deps, userID string) (follows.UserFollow, error) {
	var followData follows.UserFollow
	err := app.DB.FindOne(ctx, followingsCollection, map[string]any{"userid": userID}, &followData)
	if err != nil {
		return follows.UserFollow{}, err
	}
	return followData, nil
}

func findSuggestedUsers(ctx context.Context, app *infra.Deps, filter map[string]any) ([]UserSuggest, error) {
	var users []UserSuggest
	err := app.DB.FindMany(ctx, usersCollection, filter, &users)
	return users, err
}

func findNearbyPlaces(ctx context.Context, app *infra.Deps, filter map[string]any) ([]places.Place, error) {
	var nearbyPlaces []places.Place
	err := app.DB.FindMany(ctx, placesCollection, filter, &nearbyPlaces)
	return nearbyPlaces, err
}
