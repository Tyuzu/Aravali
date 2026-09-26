package userdata

import (
	"context"

	"scav/config"
	"scav/infra"
)

var userdataTable = config.Tables.UserDataTable

// SQLInsertUserData inserts a single user data document.
func SQLInsertUserData(ctx context.Context, app *infra.Deps, content UserData) error {
	return app.SQLDB.InsertOne(ctx, userdataTable, content)
}

// SQLDeleteUserData removes user data matching a SQL WHERE condition.
func SQLDeleteUserData(ctx context.Context, app *infra.Deps, where string, args []any) (int64, error) {
	return app.SQLDB.DeleteMany(ctx, userdataTable, where, args)
}

// SQLInsertUserDataMany inserts many user data documents.
func SQLInsertUserDataMany(ctx context.Context, app *infra.Deps, docs []any) error {
	return app.SQLDB.InsertMany(ctx, userdataTable, docs)
}

// SQLFindUserData finds user data for a given SQL WHERE query and arguments.
func SQLFindUserData(ctx context.Context, app *infra.Deps, where string, args []any, out *[]UserData) error {
	return app.SQLDB.FindMany(ctx, userdataTable, where, args, out)
}

// Database Helpers

// SQLFetchUserDataByEntity queries user data based on entity type and user ID.
func SQLFetchUserDataByEntity(ctx context.Context, app *infra.Deps, entityType, username string) ([]UserData, error) {
	where := "entity_type = $1 AND userid = $2"
	args := []any{entityType, username}

	var results []UserData
	if err := SQLFindUserData(ctx, app, where, args, &results); err != nil {
		return nil, err
	}

	if results == nil {
		return []UserData{}, nil
	}

	return results, nil
}

// SQLFetchOtherUserFeedPosts retrieves posts for a given user from the database.
func SQLFetchOtherUserFeedPosts(ctx context.Context, app *infra.Deps, username string) ([]postDoc, error) {
	where := "created_by = $1 OR username = $1"
	args := []any{username}

	var posts []postDoc
	if err := app.SQLDB.FindMany(ctx, config.Tables.FeedPostsTable, where, args, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}
