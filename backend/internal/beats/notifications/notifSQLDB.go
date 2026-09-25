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
	notifsTable      = config.Tables.NotificationsTable
	preferencesTable = config.Tables.NotificationsPreferencesTable
)

/* =========================
   DATABASE OPERATIONS
========================= */

func SQLinsertNotification(ctx context.Context, database db.Database, notif Notification) error {
	return database.Insert(ctx, notifsTable, notif)
}

func SQLinsertBulkNotifications(ctx context.Context, database db.Database, notifs []Notification) error {
	docs := make([]any, len(notifs))
	for i, v := range notifs {
		docs[i] = v
	}
	return database.InsertMany(ctx, notifsTable, docs)
}

func SQLfindNotificationsByUser(
	ctx context.Context,
	database db.Database,
	userID string,
	opts db.FindManyOptions,
	notifs *[]Notification,
) error {
	filter := map[string]any{"userid": userID}
	return database.FindManyWithOptions(ctx, notifsTable, filter, opts, notifs)
}

func SQLcountUnreadNotifications(
	ctx context.Context,
	database db.Database,
	userID string,
) (int64, error) {
	filter := map[string]any{
		"userid":  userID,
		"is_read": false,
	}
	return database.CountDocuments(ctx, notifsTable, filter)
}

func SQLupdateMarkAsRead(
	ctx context.Context,
	database db.Database,
	notificationID string,
	userID string,
) (any, error) {
	filter := map[string]any{
		"notificationid": notificationID,
		"userid":         userID,
	}
	update := map[string]any{
		"$set": map[string]any{
			"is_read": true,
		},
		"$currentDate": map[string]any{
			"updated_at": true,
		},
	}
	return database.UpdateOne(ctx, notifsTable, filter, update)
}

func SQLupdateMarkAllAsRead(
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
	return database.UpdateMany(ctx, notifsTable, filter, update)
}

func SQLdeleteNotificationByID(
	ctx context.Context,
	database db.Database,
	notificationID string,
	userID string,
) (int64, error) {
	return database.Delete(
		ctx,
		notifsTable,
		map[string]any{
			"notificationid": notificationID,
			"userid":         userID,
		},
	)
}

func SQLdeleteAllNotificationsByUser(
	ctx context.Context,
	database db.Database,
	userID string,
) error {
	return database.DeleteMany(
		ctx,
		notifsTable,
		map[string]any{"userid": userID},
	)
}

func SQLfindPreferencesByUser(
	ctx context.Context,
	database db.Database,
	userID string,
	pref *NotificationPreferences,
) error {
	return database.FindOne(
		ctx,
		preferencesTable,
		map[string]any{"userid": userID},
		pref,
	)
}

func SQLupsertPreferences(
	ctx context.Context,
	database db.Database,
	pref NotificationPreferences,
) (any, error) {
	filter := map[string]any{"userid": pref.UserID}
	update := map[string]any{"$set": pref}
	return database.UpdateOne(ctx, preferencesTable, filter, update)
}

/* =========================
   DATABASE ERROR HELPERS
========================= */

func SQLisNoDocumentsError(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}

func SQLnotificationSort() []bson.E {
	return []bson.E{
		{Key: "created_at", Value: -1},
		{Key: "notificationid", Value: -1},
	}
}
