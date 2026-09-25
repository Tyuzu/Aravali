package activity

import (
	"encoding/json"
	"net/http"
	"time"

	"scav/infra"
	"scav/utils"
)

// AnalyticsPayload represents the incoming batch payload from activityLogger.js
type AnalyticsPayload struct {
	Meta   map[string]any   `json:"meta"`
	Events []map[string]any `json:"events"`
}

// -------------------- Log Activities --------------------

func LogActivities(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var activities []Activity
		if err := json.NewDecoder(r.Body).Decode(&activities); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid payload")
			return
		}

		now := time.Now()

		for i := range activities {
			activities[i].UserID = userID
			activities[i].Timestamp = now
		}

		if err := insertActivities(r.Context(), app, activities); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to insert activities")
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, map[string]int{
			"inserted": len(activities),
		})
	}
}

// -------------------- Activity Feed --------------------

func GetActivityFeed(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		cursor, limit := parseCursor(r)

		activities, err := getActivities(
			r.Context(),
			app,
			userID,
			cursor,
			limit,
		)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "fetch failed")
			return
		}

		var nextCursor time.Time
		if len(activities) > 0 {
			nextCursor = activities[len(activities)-1].Timestamp
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"items":      activities,
			"nextCursor": nextCursor,
			"limit":      limit,
		})
	}
}

// -------------------- Analytics --------------------

func HandleAnalyticsEvent(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload AnalyticsPayload

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid payload")
			return
		}

		if len(payload.Events) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		inserted, err := insertAnalyticsEvents(
			r.Context(),
			app,
			payload,
			r.RemoteAddr,
		)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "analytics insert failed")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]int{
			"inserted": inserted,
		})
	}
}
