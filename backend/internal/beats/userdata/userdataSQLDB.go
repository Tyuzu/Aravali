// File: internal/beats/userdata/userdataSQLDB.go

package userdata

import (
	"context"

	"scav/config"
	"scav/infra"
)

var userdataTable = config.Tables.UserDataTable

// SQLInsertUserData inserts a single user data document.
func InsertUserData(ctx context.Context, app *infra.Deps, content UserData) error {

}

// SQLDeleteUserData removes user data matching a SQL WHERE condition.
func DeleteUserData(ctx context.Context, app *infra.Deps, where string, args []any) (int64, error) {

}

// SQLInsertUserDataMany inserts many user data documents.
func InsertUserDataMany(ctx context.Context, app *infra.Deps, docs []any) error {

}

// SQLFindUserData finds user data for a given SQL WHERE query and arguments.
func FindUserData(ctx context.Context, app *infra.Deps, where string, args []any, out *[]UserData) error {

}

// Database Helpers

// SQLFetchUserDataByEntity queries user data based on entity type and user ID.
func FetchUserDataByEntity(ctx context.Context, app *infra.Deps, entityType, username string) ([]UserData, error) {

	var results []UserData

	return results, nil
}

// SQLFetchOtherUserFeedPosts retrieves posts for a given user from the database.
func FetchOtherUserFeedPosts(ctx context.Context, app *infra.Deps, username string) ([]postDoc, error) {

	var posts []postDoc
	return posts, nil
}
