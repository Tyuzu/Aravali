package subscribe

import (
	"context"
	"fmt"
	"log"
	"scav/config"
	"scav/infra"
)

var subscribersTable = config.Tables.SubscribersTable
var usersTable = config.Tables.UserTable

func SQLUpdateEntitySubscription(
	ctx context.Context,
	userID,
	entityType,
	entityID,
	action string,
	app *infra.Deps,
) error {
	if action != "subscribe" && action != "unsubscribe" {
		return fmt.Errorf("invalid action: %s", action)
	}

	EnsureSubscriptionEntry(ctx, userID, app)
	EnsureSubscriptionEntry(ctx, entityID, app)

	var userUpdate any
	var entityUpdate any

	if action == "subscribe" {
		userUpdate = map[string]any{
			"$addToSet": map[string]any{"subscribed": entityID},
		}
		entityUpdate = map[string]any{
			"$addToSet": map[string]any{"subscribers": userID},
		}
	} else {
		userUpdate = map[string]any{
			"$pull": map[string]any{"subscribed": entityID},
		}
		entityUpdate = map[string]any{
			"$pull": map[string]any{"subscribers": userID},
		}
	}

	if _, err := app.DB.UpdateOne(
		ctx,
		subscribersTable,
		map[string]any{"userid": userID},
		userUpdate,
	); err != nil {
		return fmt.Errorf("failed to update user subscriptions: %w", err)
	}

	if _, err := app.DB.UpdateOne(
		ctx,
		subscribersTable,
		map[string]any{"userid": entityID},
		entityUpdate,
	); err != nil {
		return fmt.Errorf("failed to update entity subscribers: %w", err)
	}

	return nil
}

func SQLEnsureSubscriptionEntry(ctx context.Context, userID string, app *infra.Deps) {
	doc := map[string]any{
		"userid":      userID,
		"subscribed":  []string{},
		"subscribers": []string{},
	}

	err := app.DB.Upsert(
		ctx,
		subscribersTable,
		map[string]any{"userid": userID},
		map[string]any{"$setOnInsert": doc},
	)

	if err != nil {
		log.Printf("Failed to ensure subscription entry for %s: %v", userID, err)
	}
}

func SQLcountUserSubscriptionsForEntity(ctx context.Context, app *infra.Deps, userID, entityID string) (int64, error) {
	return app.DB.CountDocuments(
		ctx,
		subscribersTable,
		map[string]any{
			"userid": userID,
			"subscribed": map[string]any{
				"$in": []string{entityID},
			},
		},
	)
}

func SQLfindSubscriptionEntryByUserID(ctx context.Context, app *infra.Deps, userID string) (UserSubscribe, error) {
	var sub UserSubscribe
	err := app.DB.FindOne(ctx, subscribersTable, map[string]any{"userid": userID}, &sub)
	return sub, err
}

func SQLfindUsersByIDs(ctx context.Context, app *infra.Deps, userIDs []string) ([]map[string]any, error) {
	var subscribers []map[string]any
	err := app.DB.FindMany(
		ctx,
		usersTable,
		map[string]any{"userid": map[string]any{"$in": userIDs}},
		&subscribers,
	)
	return subscribers, err
}
