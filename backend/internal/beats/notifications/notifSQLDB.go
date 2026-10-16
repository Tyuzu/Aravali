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
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx,
		`INSERT INTO notifications (notificationid, userid, title, message, type, is_read, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) ON CONFLICT (notificationid) DO NOTHING`,
		notif.NotificationID,
		notif.UserID,
		notif.Title,
		notif.Message,
		notif.Type,
		notif.IsRead,
	)
	return err
}

func insertBulkNotifications(ctx context.Context, app *infra.Deps, notifs []Notification) error {
	if app == nil || app.SQLDB == nil || len(notifs) == 0 {
		return nil
	}
	for _, notif := range notifs {
		if err := insertNotification(ctx, app, notif); err != nil {
			return err
		}
	}
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
	if app == nil || app.SQLDB == nil || notifs == nil {
		return nil
	}
	if userID == "" {
		*notifs = []Notification{}
		return nil
	}
	query := `SELECT notificationid, userid, title, message, type, is_read, created_at, updated_at FROM notifications WHERE userid = $1 ORDER BY created_at DESC`
	limit := 20
	if v, ok := opts["limit"].(int); ok && v > 0 {
		limit = v
	}
	query += ` LIMIT $2`
	rows, err := app.SQLDB.Query(ctx, query, userID, limit)
	if err != nil {
		return err
	}
	defer rows.Close()
	items := make([]Notification, 0)
	for rows.Next() {
		var item Notification
		if err := rows.Scan(&item.NotificationID, &item.UserID, &item.Title, &item.Message, &item.Type, &item.IsRead, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return err
		}
		items = append(items, item)
	}
	*notifs = items
	return rows.Err()
}

func countUnreadNotifications(
	ctx context.Context,
	app *infra.Deps,
	userID string,
) (int64, error) {
	if app == nil || app.SQLDB == nil || userID == "" {
		return 0, nil
	}
	var count int64
	err := app.SQLDB.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE userid = $1 AND is_read = false`, userID).Scan(&count)
	return count, err
}

func updateMarkAsRead(
	ctx context.Context,
	app *infra.Deps,
	notificationID string,
	userID string,
) (int64, error) {
	if app == nil || app.SQLDB == nil || notificationID == "" {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(ctx, `UPDATE notifications SET is_read = true, updated_at = NOW() WHERE notificationid = $1 AND userid = $2`, notificationID, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func updateMarkAllAsRead(
	ctx context.Context,
	app *infra.Deps,
	userID string,
) (int64, error) {
	if app == nil || app.SQLDB == nil || userID == "" {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(ctx, `UPDATE notifications SET is_read = true, updated_at = NOW() WHERE userid = $1`, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func deleteNotificationByID(
	ctx context.Context,
	app *infra.Deps,
	notificationID string,
	userID string,
) (int64, error) {
	if app == nil || app.SQLDB == nil || notificationID == "" {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(ctx, `DELETE FROM notifications WHERE notificationid = $1 AND userid = $2`, notificationID, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func deleteAllNotificationsByUser(
	ctx context.Context,
	app *infra.Deps,
	userID string,
) (int64, error) {
	if app == nil || app.SQLDB == nil || userID == "" {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(ctx, `DELETE FROM notifications WHERE userid = $1`, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func findPreferencesByUser(
	ctx context.Context,
	app *infra.Deps,
	userID string,
	pref *NotificationPreferences,
) error {
	if app == nil || app.SQLDB == nil || pref == nil || userID == "" {
		return nil
	}
	row := app.SQLDB.QueryRow(ctx, `SELECT userid, email_notifs, push_notifs, in_app_notifs, updated_at FROM notificationpreferences WHERE userid = $1 LIMIT 1`, userID)
	return row.Scan(&pref.UserID, &pref.EmailNotifs, &pref.PushNotifs, &pref.InAppNotifs, &pref.UpdatedAt)
}

func upsertPreferences(
	ctx context.Context,
	app *infra.Deps,
	pref NotificationPreferences,
) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx,
		`INSERT INTO notificationpreferences (userid, email_notifs, push_notifs, in_app_notifs, updated_at) VALUES ($1, $2, $3, $4, NOW()) ON CONFLICT (userid) DO UPDATE SET email_notifs = EXCLUDED.email_notifs, push_notifs = EXCLUDED.push_notifs, in_app_notifs = EXCLUDED.in_app_notifs, updated_at = NOW()`,
		pref.UserID,
		pref.EmailNotifs,
		pref.PushNotifs,
		pref.InAppNotifs,
	)
	return err
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
