// File: internal/settings/settingsDB.go

package settings

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

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

func settingsJSON(settings UserSettings) map[string]any {
	return map[string]any{
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
}

func hydrateSettingsFromJSON(data map[string]any, base UserSettings) UserSettings {
	settings := base
	if data == nil {
		return settings
	}
	if value, ok := data["theme"].(string); ok && value != "" {
		settings.Theme = value
	}
	if value, ok := data["notifications"].(bool); ok {
		settings.Notifications = value
	}
	if value, ok := data["email_notifications"].(bool); ok {
		settings.EmailNotifications = value
	}
	if value, ok := data["push_notifications"].(bool); ok {
		settings.PushNotifications = value
	}
	if value, ok := data["privacy_mode"].(bool); ok {
		settings.PrivacyMode = value
	}
	if value, ok := data["profile_visibility"].(string); ok && value != "" {
		settings.ProfileVisibility = value
	}
	if value, ok := data["auto_logout"].(bool); ok {
		settings.AutoLogout = value
	}
	if value, ok := data["session_timeout"].(float64); ok {
		settings.SessionTimeout = int(value)
	} else if value, ok := data["session_timeout"].(int); ok {
		settings.SessionTimeout = value
	}
	if value, ok := data["language"].(string); ok && value != "" {
		settings.Language = value
	}
	if value, ok := data["time_zone"].(string); ok && value != "" {
		settings.TimeZone = value
	}
	if value, ok := data["currency"].(string); ok && value != "" {
		settings.Currency = value
	}
	if value, ok := data["daily_reminder"].(string); ok && value != "" {
		settings.DailyReminder = value
	}
	return settings
}

func GetUserSettingsByUserID(ctx context.Context, app *infra.Deps, userID string) (UserSettings, error) {
	if app == nil || app.SQLDB == nil {
		return DefaultSettings(userID), nil
	}
	if strings.TrimSpace(userID) == "" {
		return UserSettings{}, errors.New("user id required")
	}
	var payload []byte
	var updatedAt time.Time
	if err := app.SQLDB.QueryRow(ctx, `SELECT value, updated_at FROM settings WHERE userid = $1 LIMIT 1`, userID).Scan(&payload, &updatedAt); err != nil {
		return UserSettings{}, err
	}
	settings := DefaultSettings(userID)
	if len(payload) > 0 && string(payload) != "null" {
		var data map[string]any
		if err := json.Unmarshal(payload, &data); err == nil {
			settings = hydrateSettingsFromJSON(data, settings)
		}
	}
	settings.UserID = userID
	return settings, nil
}

func InsertUserSettings(ctx context.Context, app *infra.Deps, settings UserSettings) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if strings.TrimSpace(settings.UserID) == "" {
		return errors.New("user id required")
	}
	payload, err := json.Marshal(settingsJSON(settings))
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(ctx,
		`INSERT INTO settings (settingid, userid, key_name, value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (settingid) DO UPDATE SET
			userid = EXCLUDED.userid,
			key_name = EXCLUDED.key_name,
			value = EXCLUDED.value,
			updated_at = NOW()`,
		settings.UserID,
		settings.UserID,
		"user_settings",
		payload,
	)
	return err
}

func UpdateUserSettings(ctx context.Context, app *infra.Deps, userID string, updates map[string]any) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if strings.TrimSpace(userID) == "" {
		return errors.New("user id required")
	}
	current, err := GetUserSettingsByUserID(ctx, app, userID)
	if err != nil {
		current = DefaultSettings(userID)
	}
	mergeData := settingsJSON(current)
	for key, value := range updates {
		mergeData[key] = value
	}
	payload, err := json.Marshal(mergeData)
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(ctx,
		`UPDATE settings SET value = $1::jsonb, updated_at = NOW() WHERE userid = $2`,
		payload,
		userID,
	)
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
	if existing.UserID == "" {
		existing.UserID = settings.UserID
	}
	return sqlUpdateUserSettings(ctx, app, existing.UserID, updates)
}
