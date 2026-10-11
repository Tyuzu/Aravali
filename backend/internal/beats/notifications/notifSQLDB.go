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

func insertNotification(ctx context.Context, app *infra.Deps, notif Notification) error {
	return nil
}

func insertBulkNotifications(ctx context.Context, app *infra.Deps, notifs []Notification) error {
	return nil
}

func buildNotificationQueryOptions(page, limit int) map[string]any {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return map[string]any{
		"page":  page,
		"limit": limit,
		"sort":  "created_at DESC",
	}
}

func getUserNotificationsPage(
	ctx context.Context,
	app *infra.Deps,
	userID string,
	page, limit int,
) ([]Notification, error) {
	opts := buildNotificationQueryOptions(page, limit)
	var notifs []Notification
	if err := findNotificationsByUser(ctx, app, userID, opts, &notifs); err != nil {
		return nil, err
	}
	if notifs == nil {
		return []Notification{}, nil
	}
	return notifs, nil
}

func findNotificationsByUser(
	ctx context.Context,
	app *infra.Deps,
	userID string,
	opts map[string]any,
	notifs *[]Notification,
) error {
	if notifs != nil {
		*notifs = []Notification{}
	}
	return nil
}

func countUnreadNotifications(
	ctx context.Context,
	app *infra.Deps,
	userID string,
) (int64, error) {
	return 0, nil
}

func updateMarkAsRead(
	ctx context.Context,
	app *infra.Deps,
	notificationID string,
	userID string,
) (int64, error) {
	return 0, nil
}

func updateMarkAllAsRead(
	ctx context.Context,
	app *infra.Deps,
	userID string,
) (int64, error) {
	return 0, nil
}

func deleteNotificationByID(
	ctx context.Context,
	app *infra.Deps,
	notificationID string,
	userID string,
) (int64, error) {
	return 0, nil
}

func deleteAllNotificationsByUser(
	ctx context.Context,
	app *infra.Deps,
	userID string,
) (int64, error) {
	return 0, nil
}

func findPreferencesByUser(
	ctx context.Context,
	app *infra.Deps,
	userID string,
	pref *NotificationPreferences,
) error {
	return nil
}

func upsertPreferences(
	ctx context.Context,
	app *infra.Deps,
	pref NotificationPreferences,
) error {
	return nil
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
