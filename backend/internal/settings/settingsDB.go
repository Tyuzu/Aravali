// File: internal/settings/settingsDB.go

package settings

import (
	"context"
	"scav/config"
	"scav/infra"
)

var settingsTable = config.Tables.SettingsTable

func sqlGetUserSettingsByUserID(ctx context.Context, app *infra.Deps, userID string) (UserSettings, error) {
	return GetUserSettingsByUserID(ctx, app, userID)
}

func sqlInsertUserSettings(ctx context.Context, app *infra.Deps, settings UserSettings) error {
	return InsertUserSettings(ctx, app, settings)
}

func sqlUpdateUserSettings(ctx context.Context, app *infra.Deps, userID string, updates map[string]any) error {
	return UpdateUserSettings(ctx, app, userID, updates)
}

func sqlUpsertUserSettings(ctx context.Context, app *infra.Deps, settings UserSettings) error {
	return UpsertUserSettings(ctx, app, settings)
}

func GetUserSettingsByUserID(ctx context.Context, app *infra.Deps, userID string) (UserSettings, error) {
	var settings UserSettings
	where := "user_id = $1"
	args := []any{userID}

	if err := app.SQLDB.FindOne(ctx, settingsTable, where, args, &settings); err != nil {
		return UserSettings{}, err
	}
	return settings, nil
}

func InsertUserSettings(ctx context.Context, app *infra.Deps, settings UserSettings) error {
	return app.SQLDB.InsertOne(ctx, settingsTable, settings)
}

func UpdateUserSettings(ctx context.Context, app *infra.Deps, userID string, updates map[string]any) error {
	where := "user_id = $1"
	args := []any{userID}

	_, err := app.SQLDB.UpdateOne(ctx, settingsTable, where, args, updates)
	return err
}

func UpsertUserSettings(ctx context.Context, app *infra.Deps, settings UserSettings) error {
	existing, err := sqlGetUserSettingsByUserID(ctx, app, settings.UserID)
	if err != nil {
		return sqlInsertUserSettings(ctx, app, settings)
	}

	updates := map[string]any{
		"theme":               settings.Theme,
		"notifications":       settings.Notifications,
		"email_notifications": settings.EmailNotifications,
		"push_notifications":  settings.PushNotifications,
		"privacy_mode":        settings.PrivacyMode,
		"profile_visibility":  settings.ProfileVisibility,
		"auto_logout":         settings.AutoLogout,
		"session_timeout":     settings.SessionTimeout,
		"language":            settings.Language,
		"time_zone":           settings.TimeZone,
		"currency":            settings.Currency,
		"daily_reminder":      settings.DailyReminder,
	}

	return sqlUpdateUserSettings(ctx, app, existing.UserID, updates)
}
