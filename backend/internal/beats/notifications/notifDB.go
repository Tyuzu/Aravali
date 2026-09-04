package notifications

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"scav/config"
	db "scav/infra/db"
)

var (
	notifsCollection      = config.Collections.NotificationsCollection
	preferencesCollection = config.Collections.NotificationsPreferencesCollection
)

/* =========================
   DATABASE OPERATIONS
========================= */

func insertNotification(ctx context.Context, database db.Database, notif Notification) error {
	return database.Insert(ctx, notifsCollection, notif)
}

func insertBulkNotifications(ctx context.Context, database db.Database, notifs []Notification) error {
	docs := make([]any, len(notifs))
	for i, v := range notifs {
		docs[i] = v
	}
	return database.InsertMany(ctx, notifsCollection, docs)
}

func findNotificationsByUser(
	ctx context.Context,
	database db.Database,
	userID string,
	opts db.FindManyOptions,
	notifs *[]Notification,
) error {
	filter := map[string]any{"userid": userID}
	return database.FindManyWithOptions(ctx, notifsCollection, filter, opts, notifs)
}

func countUnreadNotifications(
	ctx context.Context,
	database db.Database,
	userID string,
) (int64, error) {
	filter := map[string]any{
		"userid":  userID,
		"is_read": false,
	}
	return database.CountDocuments(ctx, notifsCollection, filter)
}

func updateMarkAsRead(
	ctx context.Context,
	database db.Database,
	notificationID string,
) (any, error) {
	filter := map[string]any{"notificationid": notificationID}
	update := map[string]any{
		"$set": map[string]any{
			"is_read": true,
		},
		"$currentDate": map[string]any{
			"updated_at": true,
		},
	}
	return database.UpdateOne(ctx, notifsCollection, filter, update)
}

func updateMarkAllAsRead(
	ctx context.Context,
	database db.Database,
	userID string,
) (any, error) {
	filter := map[string]any{
		"userid":  userID,
		"is_read": false,
	}
	update := map[string]any{
		"$set": map[string]any{
			"is_read": true,
		},
		"$currentDate": map[string]any{
			"updated_at": true,
		},
	}
	return database.UpdateMany(ctx, notifsCollection, filter, update)
}

func deleteNotificationByID(
	ctx context.Context,
	database db.Database,
	notificationID string,
) (int64, error) {
	return database.Delete(
		ctx,
		notifsCollection,
		map[string]any{"notificationid": notificationID},
	)
}

func deleteAllNotificationsByUser(
	ctx context.Context,
	database db.Database,
	userID string,
) error {
	return database.DeleteMany(
		ctx,
		notifsCollection,
		map[string]any{"userid": userID},
	)
}

func findPreferencesByUser(
	ctx context.Context,
	database db.Database,
	userID string,
	pref *NotificationPreferences,
) error {
	return database.FindOne(
		ctx,
		preferencesCollection,
		map[string]any{"userid": userID},
		pref,
	)
}

func upsertPreferences(
	ctx context.Context,
	database db.Database,
	pref NotificationPreferences,
) (any, error) {
	filter := map[string]any{"userid": pref.UserID}
	update := map[string]any{"$set": pref}
	return database.UpdateOne(ctx, preferencesCollection, filter, update)
}

/* =========================
   DATABASE ERROR HELPERS
========================= */

func isNoDocumentsError(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}

func notificationSort() []bson.E {
	return []bson.E{
		{Key: "created_at", Value: -1},
		{Key: "notificationid", Value: -1},
	}
}
