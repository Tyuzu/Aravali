package notifications

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"scav/config"
	"scav/infra/sqldb"
)

var (
	notifsTable      = config.Tables.NotificationsTable
	preferencesTable = config.Tables.NotificationsPreferencesTable
)

/* =========================
   DATABASE OPERATIONS
========================= */

func SQLinsertNotification(ctx context.Context, database sqldb.Database, notif Notification) error {
	return database.InsertOne(ctx, notifsTable, notif)
}

func SQLinsertBulkNotifications(ctx context.Context, database sqldb.Database, notifs []Notification) error {
	docs := make([]any, len(notifs))
	for i, v := range notifs {
		docs[i] = v
	}
	return database.InsertMany(ctx, notifsTable, docs)
}

func SQLfindNotificationsByUser(
	ctx context.Context,
	database sqldb.Database,
	userID string,
	opts sqldb.FindManyOptions,
	notifs *[]Notification,
) error {
	where := "userid = $1"
	args := []any{userID}

	return database.FindManyWithOptions(ctx, notifsTable, where, args, opts, notifs)
}

func SQLcountUnreadNotifications(
	ctx context.Context,
	database sqldb.Database,
	userID string,
) (int64, error) {
	where := "userid = $1 AND is_read = false"
	args := []any{userID}

	return database.Count(ctx, notifsTable, where, args)
}

func SQLupdateMarkAsRead(
	ctx context.Context,
	database sqldb.Database,
	notificationID string,
	userID string,
) (int64, error) {
	where := "notificationid = $1 AND userid = $2"
	args := []any{notificationID, userID}

	updateValues := map[string]any{
		"is_read":    true,
		"updated_at": time.Now(),
	}

	return database.UpdateOne(ctx, notifsTable, where, args, updateValues)
}

func SQLupdateMarkAllAsRead(
	ctx context.Context,
	database sqldb.Database,
	userID string,
) (int64, error) {
	where := "userid = $1 AND is_read = false"
	args := []any{userID}

	updateValues := map[string]any{
		"is_read":    true,
		"updated_at": time.Now(),
	}

	return database.UpdateMany(ctx, notifsTable, where, args, updateValues)
}

func SQLdeleteNotificationByID(
	ctx context.Context,
	database sqldb.Database,
	notificationID string,
	userID string,
) (int64, error) {
	where := "notificationid = $1 AND userid = $2"
	args := []any{notificationID, userID}

	return database.DeleteOne(ctx, notifsTable, where, args)
}

func SQLdeleteAllNotificationsByUser(
	ctx context.Context,
	database sqldb.Database,
	userID string,
) (int64, error) {
	where := "userid = $1"
	args := []any{userID}

	return database.DeleteMany(ctx, notifsTable, where, args)
}

func SQLfindPreferencesByUser(
	ctx context.Context,
	database sqldb.Database,
	userID string,
	pref *NotificationPreferences,
) error {
	where := "userid = $1"
	args := []any{userID}

	return database.FindOne(ctx, preferencesTable, where, args, pref)
}

func SQLupsertPreferences(
	ctx context.Context,
	database sqldb.Database,
	pref NotificationPreferences,
) error {
	return database.Upsert(ctx, preferencesTable, "userid", pref)
}

/* =========================
   DATABASE ERROR HELPERS
========================= */

func SQLisNoDocumentsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func SQLnotificationSort() string {
	return "created_at DESC, notificationid DESC"
}
