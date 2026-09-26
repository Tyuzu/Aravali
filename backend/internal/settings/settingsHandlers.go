package settings

import (
	"encoding/json"
	"net/http"

	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	"scav/utils"
)

func GetSettings(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := getUserID(r)
		if !ok {
			utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "unauthorized",
			})
			return
		}

		settings, err := sqlGetUserSettingsByUserID(ctx, app, userID)
		if err != nil {
			settings = DefaultSettings(userID)
			_ = sqlInsertUserSettings(ctx, app, settings)
		}

		utils.RespondWithJSON(w, http.StatusOK, settings)
	}
}

func GetSettingsSchema(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, http.StatusOK, settingsSchema)
	}
}

func UpdateSettings(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := getUserID(r)
		if !ok {
			utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "unauthorized",
			})
			return
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid payload",
			})
			return
		}

		allowed := map[string]bool{
			"theme":               true,
			"notifications":       true,
			"email_notifications": true,
			"push_notifications":  true,
			"privacy_mode":        true,
			"profile_visibility":  true,
			"auto_logout":         true,
			"session_timeout":     true,
			"language":            true,
			"time_zone":           true,
			"currency":            true,
			"daily_reminder":      true,
		}

		updateFields := map[string]any{}

		for key, value := range payload {
			if !allowed[key] {
				utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{
					"error": "invalid setting type: " + key,
				})
				return
			}

			if err := validateSetting(key, value); err != nil {
				utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{
					"error": err.Error(),
				})
				return
			}

			updateFields[key] = value
		}

		if len(updateFields) == 0 {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{
				"error": "empty update payload",
			})
			return
		}

		if err := sqlUpdateUserSettings(ctx, app, userID, updateFields); err != nil {
			settings := DefaultSettings(userID)
			for k, v := range updateFields {
				applyPatch(&settings, k, v)
			}
			_ = sqlInsertUserSettings(ctx, app, settings)
		}

		mqpayload, _ := json.Marshal(mqevent.UserSettingsUpdatedPayload{})
		mq.PublishWithMeta(ctx, app.MQ, mqevent.UserSettingsUpdatedEvent, mqpayload)

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"status":  "success",
			"message": "settings updated",
			"data":    updateFields,
		})
	}
}

func ResetSettings(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := getUserID(r)
		if !ok {
			utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "unauthorized",
			})
			return
		}

		defaults := DefaultSettings(userID)
		if err := sqlUpsertUserSettings(ctx, app, defaults); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to reset settings",
			})
			return
		}

		mqpayload, _ := json.Marshal(mqevent.UserSettingsResetPayload{})
		mq.PublishWithMeta(ctx, app.MQ, mqevent.UserSettingsResetEvent, mqpayload)

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"status":  "success",
			"message": "settings reset to defaults",
		})
	}
}

func InitUserSettings(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := getUserID(r)
		if !ok {
			utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "unauthorized",
			})
			return
		}

		_, err := sqlGetUserSettingsByUserID(ctx, app, userID)
		if err == nil {
			utils.RespondWithJSON(w, http.StatusOK, false)
			return
		}

		defaults := DefaultSettings(userID)
		if err := sqlInsertUserSettings(ctx, app, defaults); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to initialize settings",
			})
			return
		}

		mqpayload, _ := json.Marshal(mqevent.UserSettingsInitiatedPayload{})
		mq.PublishWithMeta(ctx, app.MQ, mqevent.UserSettingsInitiatedEvent, mqpayload)

		utils.RespondWithJSON(w, http.StatusOK, true)
	}
}
