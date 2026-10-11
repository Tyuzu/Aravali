// File: internal/beats/follows/followsSQLDB.go

package follows

import (
	"context"
	"errors"
	"time"

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
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if currentUserID == "" || targetUserID == "" || currentUserID == targetUserID {
		return errors.New("invalid follow request")
	}

	switch action {
	case "follow":
		_, err := app.SQLDB.Exec(ctx,
			`INSERT INTO `+followingsTable+` (userid, followingid, created_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (userid, followingid) DO NOTHING`,
			currentUserID, targetUserID, time.Now(),
		)
		return err
	case "unfollow":
		_, err := app.SQLDB.Exec(ctx,
			`DELETE FROM `+followingsTable+` WHERE userid = $1 AND followingid = $2`,
			currentUserID, targetUserID,
		)
		return err
	default:
		return errors.New("unsupported follow action")
	}
}

func CreateFollowEntry(userid string, app *infra.Deps) {
	if app == nil || app.SQLDB == nil || userid == "" {
		return
	}
	_, _ = app.SQLDB.Exec(context.Background(),
		`INSERT INTO `+followingsTable+` (userid, followingid, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING`,
		userid, userid, time.Now(),
	)
}

func CountFollowRelationship(ctx context.Context, app *infra.Deps, userID, followedUserID string) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	var count int64
	err := app.SQLDB.QueryRow(ctx,
		`SELECT COUNT(*) FROM `+followingsTable+` WHERE userid = $1 AND followingid = $2`,
		userID, followedUserID,
	).Scan(&count)
	return count, err
}

func FindFollowEntryByUserID(ctx context.Context, app *infra.Deps, userID string, out *UserFollow) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	var foundUserID string
	if err := app.SQLDB.QueryRow(ctx, `SELECT userid FROM `+followingsTable+` WHERE userid = $1 LIMIT 1`, userID).Scan(&foundUserID); err != nil {
		return err
	}
	out.UserID = foundUserID
	return nil
}

func FindUsersByIDsForFollow(ctx context.Context, app *infra.Deps, userIDs []string, out *[]auth.User) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if len(userIDs) == 0 {
		*out = []auth.User{}
		return nil
	}
	rows, err := app.SQLDB.Query(ctx, `SELECT * FROM `+usersTable+` WHERE userid = ANY($1)`, userIDs)
	if err != nil {
		return err
	}
	defer rows.Close()
	results := make([]auth.User, 0)
	for rows.Next() {
		var user auth.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Name); err == nil {
			results = append(results, user)
		}
	}
	*out = results
	return rows.Err()
}

// GetUserFollowData returns followers and follows for a user.
func GetUserFollowData(ctx context.Context, userID string, app infra.Deps) (UserFollow, error) {
	var uf UserFollow
	if userID == "" {
		return uf, nil
	}
	if app.SQLDB == nil {
		return uf, nil
	}
	uf.UserID = userID
	return uf, nil
}
