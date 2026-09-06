package follows

import (
	"context"
	"fmt"
	"scav/config"
	"scav/infra"
	"scav/internal/auth"
	log "scav/utils/logger"
	"time"
)

var followingsCollection = config.Collections.FollowingsCollection
var usersCollection = config.Collections.UserCollection

func UpdateFollowRelationship(
	ctx context.Context,
	currentUserID,
	targetUserID,
	action string,
	app *infra.Deps,
) error {
	if action != "follow" && action != "unfollow" {
		return fmt.Errorf("invalid action: %s", action)
	}

	var currentUserUpdate any
	var targetUserUpdate any

	if action == "follow" {
		currentUserUpdate = map[string]any{
			"$addToSet": map[string]any{"follows": targetUserID},
		}
		targetUserUpdate = map[string]any{
			"$addToSet": map[string]any{"followers": currentUserID},
		}
	} else {
		currentUserUpdate = map[string]any{
			"$pull": map[string]any{"follows": targetUserID},
		}
		targetUserUpdate = map[string]any{
			"$pull": map[string]any{"followers": currentUserID},
		}
	}

	if err := app.DB.Upsert(
		ctx,
		followingsCollection,
		map[string]any{"userid": currentUserID},
		currentUserUpdate,
	); err != nil {
		return fmt.Errorf("failed to update current user's follows: %w", err)
	}

	if err := app.DB.Upsert(
		ctx,
		followingsCollection,
		map[string]any{"userid": targetUserID},
		targetUserUpdate,
	); err != nil {
		return fmt.Errorf("failed to update target user's followers: %w", err)
	}

	return nil
}

func CreateFollowEntry(userid string, app *infra.Deps) {
	follow := UserFollow{
		UserID:    userid,
		Follows:   []string{},
		Followers: []string{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := app.DB.Insert(
		ctx,
		followingsCollection,
		follow,
	)
	if err != nil {
		log.Printf("Error inserting follow entry for %s: %v", userid, err)
	}
}

func CountFollowRelationship(ctx context.Context, app *infra.Deps, userID, followedUserID string) (int64, error) {
	return app.DB.CountDocuments(
		ctx,
		followingsCollection,
		map[string]any{
			"userid": userID,
			"follows": map[string]any{
				"$in": []string{followedUserID},
			},
		},
	)
}

func FindFollowEntryByUserID(ctx context.Context, app *infra.Deps, userID string, out *UserFollow) error {
	return app.DB.FindOne(
		ctx,
		followingsCollection,
		map[string]any{"userid": userID},
		out,
	)
}

func FindUsersByIDsForFollow(ctx context.Context, app *infra.Deps, userIDs []string, out *[]auth.User) error {
	return app.DB.FindMany(
		ctx,
		usersCollection,
		map[string]any{
			"userid": map[string]any{"$in": userIDs},
		},
		out,
	)
}
