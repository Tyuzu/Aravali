package places

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"scav/config"
	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	placedb "scav/internal/places/placedb"
	"scav/utils"
)

// --- EditPlace endpoint ---
func EditPlace(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		placeID := strings.TrimSpace(utils.GetParam(r, "placeid"))
		if placeID == "" {
			http.Error(w, "Place ID is required", http.StatusBadRequest)
			return
		}

		requestingUserID, ok := ctx.Value(config.UserIDKey).(string)
		if !ok {
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		}

		// Fetch existing place (use placeid)
		var existing struct {
			CreatedBy string `bson:"createdBy"`
		}
		if err := placedb.FindOnePlace(
			ctx,
			app,
			map[string]any{"placeid": placeID},
			&existing,
		); err != nil {
			http.Error(w, "Place not found", http.StatusNotFound)
			return
		}

		if existing.CreatedBy != requestingUserID {
			http.Error(w, "You are not authorized to edit this place", http.StatusForbidden)
			return
		}

		// Parse update fields
		_, updateFields, err := parseAndBuildPlace(r, "update")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(updateFields) == 0 {
			http.Error(w, "No fields to update", http.StatusBadRequest)
			return
		}

		updateFields["updated_at"] = time.Now()
		updateFields["updatedBy"] = requestingUserID

		// ✅ Update using placeid and plain fields
		if _, err := placedb.UpdatePlace(
			ctx,
			app,
			map[string]any{"placeid": placeID},
			updateFields,
		); err != nil {
			http.Error(w, "Failed to update place", http.StatusInternalServerError)
			return
		}

		mqpayload, _ := json.Marshal(mqevent.PlaceUpdatedPayload{})

		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.PlaceUpdatedEvent, mqpayload)

		utils.RespondWithJSON(w, http.StatusOK, updateFields)
	}
}

// --- DeletePlace endpoint ---
func DeletePlace(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		placeID := strings.TrimSpace(utils.GetParam(r, "placeid"))
		if placeID == "" {
			http.Error(w, "Place ID required", http.StatusBadRequest)
			return
		}

		userID, ok := ctx.Value(config.UserIDKey).(string)
		if !ok {
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		}

		var existing struct{ CreatedBy string `bson:"createdBy"` }
		if err := placedb.FindOnePlace(ctx, app, map[string]any{"placeid": placeID}, &existing); err != nil {
			http.Error(w, "Place not found", http.StatusNotFound)
			return
		}

		if existing.CreatedBy != userID {
			http.Error(w, "Not authorized", http.StatusForbidden)
			return
		}

		if _, err := placedb.DeletePlace(ctx, app, map[string]any{"placeid": placeID}); err != nil {
			http.Error(w, "Failed to delete place", http.StatusInternalServerError)
			return
		}

		mqpayload, _ := json.Marshal(mqevent.PlaceDeletedPayload{})
		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.PlaceDeletedEvent, mqpayload)

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"placeid": placeID, "deleted": true})
	}
}
