package follows

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"scav/config"
	"scav/infra"
	"scav/internal/auth"
	log "scav/utils/logger"
)

var followingsTable = config.Tables.FollowingsTable
var usersTable = config.Tables.UserTable

func SQLUpdateFollowRelationship(
	ctx context.Context,
	currentUserID,
	targetUserID,
	action string,
	app *infra.Deps,
) error {
	if action != "follow" && action != "unfollow" {
		return fmt.Errorf("invalid action: %s", action)
	}

	return app.SQLDB.RunTransaction(ctx, func(tx pgx.Tx) error {
		// Ensure initial follow entries exist for both users
		if err := app.SQLDB.Upsert(ctx, followingsTable, "userid", UserFollow{
			UserID:    currentUserID,
			Follows:   []string{},
			Followers: []string{},
		}); err != nil {
			return fmt.Errorf("failed to ensure current user follow entry: %w", err)
		}

		if err := app.SQLDB.Upsert(ctx, followingsTable, "userid", UserFollow{
			UserID:    targetUserID,
			Follows:   []string{},
			Followers: []string{},
		}); err != nil {
			return fmt.Errorf("failed to ensure target user follow entry: %w", err)
		}

		// Update follows for current user
		var currentUserQuery string
		if action == "follow" {
			currentUserQuery = fmt.Sprintf(`
				UPDATE %s 
				SET follows = array_append(follows, $1) 
				WHERE userid = $2 AND NOT ($1 = ANY(follows))`, followingsTable)
		} else {
			currentUserQuery = fmt.Sprintf(`
				UPDATE %s 
				SET follows = array_remove(follows, $1) 
				WHERE userid = $2`, followingsTable)
		}

		if _, err := tx.Exec(ctx, currentUserQuery, targetUserID, currentUserID); err != nil {
			return fmt.Errorf("failed to update current user's follows: %w", err)
		}

		// Update followers for target user
		var targetUserQuery string
		if action == "follow" {
			targetUserQuery = fmt.Sprintf(`
				UPDATE %s 
				SET followers = array_append(followers, $1) 
				WHERE userid = $2 AND NOT ($1 = ANY(followers))`, followingsTable)
		} else {
			targetUserQuery = fmt.Sprintf(`
				UPDATE %s 
				SET followers = array_remove(followers, $1) 
				WHERE userid = $2`, followingsTable)
		}

		if _, err := tx.Exec(ctx, targetUserQuery, currentUserID, targetUserID); err != nil {
			return fmt.Errorf("failed to update target user's followers: %w", err)
		}

		return nil
	})
}

func SQLCreateFollowEntry(userid string, app *infra.Deps) {
	follow := UserFollow{
		UserID:    userid,
		Follows:   []string{},
		Followers: []string{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := app.SQLDB.InsertOne(ctx, followingsTable, follow); err != nil {
		log.Printf("Error inserting follow entry for %s: %v", userid, err)
	}
}

func SQLCountFollowRelationship(ctx context.Context, app *infra.Deps, userID, followedUserID string) (int64, error) {
	where := "userid = $1 AND $2 = ANY(follows)"
	args := []any{userID, followedUserID}

	return app.SQLDB.Count(ctx, followingsTable, where, args)
}

func SQLFindFollowEntryByUserID(ctx context.Context, app *infra.Deps, userID string, out *UserFollow) error {
	where := "userid = $1"
	args := []any{userID}

	return app.SQLDB.FindOne(ctx, followingsTable, where, args, out)
}

func SQLFindUsersByIDsForFollow(ctx context.Context, app *infra.Deps, userIDs []string, out *[]auth.User) error {
	if len(userIDs) == 0 {
		return nil
	}

	where := "userid = ANY($1)"
	args := []any{userIDs}

	return app.SQLDB.FindMany(ctx, usersTable, where, args, out)
}
