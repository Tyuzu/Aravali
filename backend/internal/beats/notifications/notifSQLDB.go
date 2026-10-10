// File: internal/beats/notifications/notifSQLDB.go

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

func insertNotification(ctx context.Context, database sqldb.Database, notif Notification) error {
	return database.InsertOne(ctx, notifsTable, notif)
}

func insertBulkNotifications(ctx context.Context, database sqldb.Database, notifs []Notification) error {
	docs := make([]any, len(notifs))
	for i, v := range notifs {
		docs[i] = v
	}
	return database.InsertMany(ctx, notifsTable, docs)
}

func buildNotificationQueryOptions(page, limit int) sqldb.FindManyOptions {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	skip := (page - 1) * limit
	return sqldb.FindManyOptions{
		Limit:   int64(limit),
		Offset:  int64(skip),
		OrderBy: notificationSort(),
	}
}

func getUserNotificationsPage(
	ctx context.Context,
	database sqldb.Database,
	userID string,
	page, limit int,
) ([]Notification, error) {
	opts := buildNotificationQueryOptions(page, limit)
	var notifs []Notification
	if err := findNotificationsByUser(ctx, database, userID, opts, &notifs); err != nil {
		return nil, err
	}
	if notifs == nil {
		return []Notification{}, nil
	}
	return notifs, nil
}

func findNotificationsByUser(
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

func countUnreadNotifications(
	ctx context.Context,
	database sqldb.Database,
	userID string,
) (int64, error) {
	where := "userid = $1 AND is_read = false"
	args := []any{userID}

	return database.Count(ctx, notifsTable, where, args)
}

func updateMarkAsRead(
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

func updateMarkAllAsRead(
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

func deleteNotificationByID(
	ctx context.Context,
	database sqldb.Database,
	notificationID string,
	userID string,
) (int64, error) {
	where := "notificationid = $1 AND userid = $2"
	args := []any{notificationID, userID}

	return database.DeleteOne(ctx, notifsTable, where, args)
}

func deleteAllNotificationsByUser(
	ctx context.Context,
	database sqldb.Database,
	userID string,
) (int64, error) {
	where := "userid = $1"
	args := []any{userID}

	return database.DeleteMany(ctx, notifsTable, where, args)
}

func findPreferencesByUser(
	ctx context.Context,
	database sqldb.Database,
	userID string,
	pref *NotificationPreferences,
) error {
	where := "userid = $1"
	args := []any{userID}

	return database.FindOne(ctx, preferencesTable, where, args, pref)
}

func upsertPreferences(
	ctx context.Context,
	database sqldb.Database,
	pref NotificationPreferences,
) error {
	return database.Upsert(ctx, preferencesTable, "userid", pref)
}

/* =========================
   DATABASE ERROR HELPERS
========================= */

func isNoDocumentsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func notificationSort() string {
	return "created_at DESC, notificationid DESC"
}
