// File: internal/beats/follows/followsSQLDB.go

package follows

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/internal/auth"
)

var followingsTable = config.Tables.FollowingsTable
var usersTable = config.Tables.UserTable

func UpdateFollowRelationship(
	ctx context.Context,
	currentUserID,
	targetUserID,
	action string,
	app *infra.Deps,
) error {
}

func CreateFollowEntry(userid string, app *infra.Deps) {

}

func CountFollowRelationship(ctx context.Context, app *infra.Deps, userID, followedUserID string) (int64, error) {
}

func FindFollowEntryByUserID(ctx context.Context, app *infra.Deps, userID string, out *UserFollow) error {

}

func FindUsersByIDsForFollow(ctx context.Context, app *infra.Deps, userIDs []string, out *[]auth.User) error {

}

// GetUserFollowData returns followers and follows for a user.
func GetUserFollowData(ctx context.Context, userID string, app infra.Deps) (UserFollow, error) {
	var uf UserFollow

	return uf, nil
}
