package admin

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"scav/config"
	"scav/config/mqevent"
	"scav/infra"
	"scav/infra/mq"
	"scav/utils"

	"go.mongodb.org/mongo-driver/bson"
)

var moderatorApplicationsCollection = config.Collections.ModeratorApplications

type ModeratorApplication struct {
	ID        string    `json:"id" bson:"id"`
	UserID    string    `json:"userid" bson:"userid"`
	Reason    string    `json:"reason" bson:"reason"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func ApplyModerator(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var payload struct {
			UserID string `json:"userid"`
			Reason string `json:"reason"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, `{"error":"Invalid JSON payload"}`, http.StatusBadRequest)
			return
		}

		payload.UserID = strings.TrimSpace(payload.UserID)
		payload.Reason = strings.TrimSpace(payload.Reason)

		if payload.UserID == "" || payload.Reason == "" {
			http.Error(w, `{"error":"Missing required fields"}`, http.StatusBadRequest)
			return
		}

		var existing ModeratorApplication
		err := app.DB.FindOne(
			ctx,
			moderatorApplicationsCollection,
			bson.M{"userid": payload.UserID},
			&existing,
		)
		if err == nil {
			http.Error(w, `{"error":"You have already applied to be a moderator"}`, http.StatusConflict)
			return
		}

		now := time.Now().UTC()
		appx := ModeratorApplication{
			ID:        "mod_" + utils.GenerateRandomString(16),
			UserID:    payload.UserID,
			Reason:    payload.Reason,
			Status:    "pending",
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := app.DB.Insert(ctx, moderatorApplicationsCollection, appx); err != nil {
			http.Error(w, `{"error":"Failed to save application"}`, http.StatusInternalServerError)
			return
		}

		mqpayload, _ := json.Marshal(mqevent.AppliedForModeratorRolePayload{})
		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.AppliedForModeratorRoleEvent, mqpayload)

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message": "Moderator application submitted",
			"id":      appx.ID,
		})
	}
}

func ListModeratorApplications(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		status := r.URL.Query().Get("status")

		filter := bson.M{}
		if status != "" {
			filter["status"] = status
		}

		var applications []ModeratorApplication
		err := app.DB.FindMany(ctx, moderatorApplicationsCollection, filter, &applications)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch applications",
			})
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, applications)
	}
}

func ApproveModerator(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := utils.GetParam(r, "id")
		if id == "" {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid or missing application ID",
			})
			return
		}

		now := time.Now().UTC()
		_, err := app.DB.UpdateOne(
			ctx,
			moderatorApplicationsCollection,
			bson.M{"id": id},
			bson.M{
				"$set": bson.M{
					"status":     "approved",
					"updatedAt":  now,
					"updated_at": now,
				},
			},
		)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{
				"error": "Application not found or update failed",
			})
			return
		}

		mqpayload, _ := json.Marshal(mqevent.ApprovedModeratorRoleRequestPayload{
			ApplicationID: id,
			ApprovedAt:    now,
		})
		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.ApprovedModeratorRoleRequestEvent, mqpayload)

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Application approved successfully",
		})
	}
}

func RejectModerator(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := utils.GetParam(r, "id")
		if id == "" {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid or missing application ID",
			})
			return
		}

		now := time.Now().UTC()
		_, err := app.DB.UpdateOne(
			ctx,
			moderatorApplicationsCollection,
			bson.M{"id": id},
			bson.M{
				"$set": bson.M{
					"status":     "rejected",
					"updatedAt":  now,
					"updated_at": now,
				},
			},
		)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{
				"error": "Application not found or update failed",
			})
			return
		}

		mqpayload, _ := json.Marshal(mqevent.RejectedModeratorRoleRequestPayload{
			ApplicationID: id,
			RejectedAt:    now,
		})
		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.RejectedModeratorRoleRequestEvent, mqpayload)

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Application rejected successfully",
		})
	}
}
