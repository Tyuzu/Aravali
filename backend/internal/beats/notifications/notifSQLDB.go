// File: internal/beats/notifications/notifSQLDB.go

package notifications

import (
	"context"
	"database/sql"
	"errors"

	"scav/config"
	"scav/infra"
)

var (
	notifsTable      = config.Tables.NotificationsTable
	preferencesTable = config.Tables.NotificationsPreferencesTable
)

/* =========================
   DATABASE OPERATIONS
========================= */

func insertNotification(ctx context.Context, app infra.Deps, notif Notification) error {

}

func insertBulkNotifications(ctx context.Context, app infra.Deps, notifs []Notification) error {

}

func buildNotificationQueryOptions(page, limit int) map[string]any {

}

func getUserNotificationsPage(
	ctx context.Context,
	app infra.Deps,
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
	app infra.Deps,
	userID string,
	opts map[string]any,
	notifs *[]Notification,
) error {

}

func countUnreadNotifications(
	ctx context.Context,
	app infra.Deps,
	userID string,
) (int64, error) {

}

func updateMarkAsRead(
	ctx context.Context,
	app infra.Deps,
	notificationID string,
	userID string,
) (int64, error) {

}

func updateMarkAllAsRead(
	ctx context.Context,
	app infra.Deps,
	userID string,
) (int64, error) {

}

func deleteNotificationByID(
	ctx context.Context,
	app infra.Deps,
	notificationID string,
	userID string,
) (int64, error) {

}

func deleteAllNotificationsByUser(
	ctx context.Context,
	app infra.Deps,
	userID string,
) (int64, error) {

}

func findPreferencesByUser(
	ctx context.Context,
	app infra.Deps,
	userID string,
	pref *NotificationPreferences,
) error {

}

func upsertPreferences(
	ctx context.Context,
	app infra.Deps,
	pref NotificationPreferences,
) error {

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
