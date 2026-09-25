package userdata

import (
	"context"
	"scav/config"
	"scav/infra"
)

var userdataTable = config.Tables.UserDataTable

// InsertUserData inserts a single user data document.
func SQLInsertUserData(ctx context.Context, app *infra.Deps, content UserData) error {
	return app.DB.InsertOne(ctx, userdataTable, content)
}

// DeleteUserData removes user data matching filter.
func SQLDeleteUserData(ctx context.Context, app *infra.Deps, filter map[string]any) error {
	return app.DB.DeleteMany(ctx, userdataTable, filter)
}

// InsertUserDataMany inserts many user data documents.
func SQLInsertUserDataMany(ctx context.Context, app *infra.Deps, docs []any) error {
	return app.DB.InsertMany(ctx, userdataTable, docs)
}

// FindUserData finds user data for a given filter.
func SQLFindUserData(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]UserData) error {
	return app.DB.FindMany(ctx, userdataTable, filter, out)
}

// Database Helpers

// FetchUserDataByEntity queries user data based on entity type and user ID.
func SQLFetchUserDataByEntity(ctx context.Context, app *infra.Deps, entityType, username string) ([]UserData, error) {
	filter := map[string]any{
		"entity_type": entityType,
		"userid":      username,
	}

	var results []UserData
	if err := FindUserData(ctx, app, filter, &results); err != nil {
		return nil, err
	}

	if results == nil {
		return []UserData{}, nil
	}

	return results, nil
}

// FetchOtherUserFeedPosts retrieves posts for a given user from the database.
func SQLFetchOtherUserFeedPosts(ctx context.Context, app *infra.Deps, username string) ([]postDoc, error) {
	filter := map[string]any{
		"$or": []map[string]any{
			{"createdBy": username},
			{"username": username},
		},
	}

	var posts []postDoc
	if err := app.DB.FindMany(ctx, config.Tables.FeedPostsTable, filter, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}
