package subscribe

import (
	"context"
	"fmt"
	"log"
	"scav/config"
	"scav/infra"
)

var subscribersCollection = config.Collections.SubscribersCollection
var usersCollection = config.Collections.UserCollection

func UpdateEntitySubscription(
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
		subscribersCollection,
		map[string]any{"userid": userID},
		userUpdate,
	); err != nil {
		return fmt.Errorf("failed to update user subscriptions: %w", err)
	}

	if _, err := app.DB.UpdateOne(
		ctx,
		subscribersCollection,
		map[string]any{"userid": entityID},
		entityUpdate,
	); err != nil {
		return fmt.Errorf("failed to update entity subscribers: %w", err)
	}

	return nil
}

func EnsureSubscriptionEntry(ctx context.Context, userID string, app *infra.Deps) {
	doc := map[string]any{
		"userid":      userID,
		"subscribed":  []string{},
		"subscribers": []string{},
	}

	err := app.DB.Upsert(
		ctx,
		subscribersCollection,
		map[string]any{"userid": userID},
		map[string]any{"$setOnInsert": doc},
	)

	if err != nil {
		log.Printf("Failed to ensure subscription entry for %s: %v", userID, err)
	}
}
