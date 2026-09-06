package merch

import (
	"encoding/json"
	"net/http"
	"scav/infra/mq"
	"time"

	"scav/config"
	"scav/config/mqevent"
	"scav/infra"
	"scav/internal/beats/auditlog"
	"scav/utils"
)

func validateEntityType(t string) bool {
	return t == "event" || t == "farm" || t == "artist"
}

// ---------------------- Create Merch ----------------------
func CreateMerch(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		entityType := utils.GetParam(r, "entityType")
		eventID := utils.GetParam(r, "eventid")

		if !validateEntityType(entityType) {
			utils.RespondWithJSON(w, 400, map[string]any{"success": false, "error": "invalid entity type"})
			return
		}

		// SECURITY: Verify user is authenticated
		userID, ok := r.Context().Value(config.UserIDKey).(string)
		if !ok || userID == "" {
			utils.RespondWithJSON(w, 401, map[string]any{"success": false, "error": "unauthorized"})
			return
		}

		// SECURITY: Verify user is the owner of the entity
		if _, err := getEntityOwner(r.Context(), app, entityType, eventID); err != nil {
			if err.Error() == "cannot verify ownership" {
				utils.RespondWithJSON(w, 403, map[string]any{"success": false, "error": "cannot verify ownership"})
				return
			}
			utils.RespondWithJSON(w, 404, map[string]any{"success": false, "error": "entity not found"})
			return
		}

		if owner, err := getEntityOwner(r.Context(), app, entityType, eventID); err != nil {
			utils.RespondWithJSON(w, 403, map[string]any{"success": false, "error": "forbidden: only entity owner can create merch"})
			return
		} else if owner != userID {
			utils.RespondWithJSON(w, 403, map[string]any{"success": false, "error": "forbidden: only entity owner can create merch"})
			return
		}

		var body struct {
			Name     string  `json:"name"`
			Price    float64 `json:"price"`
			Discount float64 `json:"discount"`
			Stock    int     `json:"stock"`
			Photo    string  `json:"merch_pic"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.RespondWithJSON(w, 400, map[string]any{"success": false, "error": "invalid json"})
			return
		}

		if body.Name == "" || body.Price <= 0 || body.Stock < 0 {
			utils.RespondWithJSON(w, 400, map[string]any{"success": false, "error": "invalid merch data"})
			return
		}

		now := time.Now()
		merch := Merch{
			MerchID:    utils.GenerateRandomString(14),
			EntityType: entityType,
			EntityID:   eventID,
			Name:       body.Name,
			Price:      body.Price,
			Discount:   body.Discount,
			Stock:      body.Stock,
			MerchPhoto: body.Photo,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		if err := insertMerch(r.Context(), app, merch); err != nil {
			utils.RespondWithJSON(w, 500, map[string]any{"success": false, "error": "insert failed"})
			return
		}

		// Log audit trail for merchandise creation
		auditlog.LogAction(
			r.Context(), app, r, userID,
			auditlog.AuditActionMerchCreate,
			"merchandise", merch.MerchID, "success",
			map[string]interface{}{
				"name":        merch.Name,
				"price":       merch.Price,
				"stock":       merch.Stock,
				"entity_type": entityType,
				"entity_id":   eventID,
			},
		)

		mqpayload, _ := json.Marshal(mqevent.MerchCreatedPayload{})

		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.MerchCreatedEvent, mqpayload)

		utils.RespondWithJSON(w, 201, map[string]any{"success": true, "data": merch})
	}
}

// ---------------------- Edit Merch ----------------------
func EditMerch(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		entityType := utils.GetParam(r, "entityType")
		eventID := utils.GetParam(r, "eventid")
		merchID := utils.GetParam(r, "merchid")

		// SECURITY: Verify user is authenticated
		userID, ok := r.Context().Value(config.UserIDKey).(string)
		if !ok || userID == "" {
			utils.RespondWithJSON(w, 401, map[string]any{"success": false, "error": "unauthorized"})
			return
		}

		// SECURITY: Verify user is the owner of the entity
		if owner, err := getEntityOwner(r.Context(), app, entityType, eventID); err != nil {
			if err.Error() == "cannot verify ownership" {
				utils.RespondWithJSON(w, 403, map[string]any{"success": false, "error": "cannot verify ownership"})
				return
			}
			utils.RespondWithJSON(w, 404, map[string]any{"success": false, "error": "entity not found"})
			return
		} else if owner != userID {
			utils.RespondWithJSON(w, 403, map[string]any{"success": false, "error": "forbidden: only entity owner can edit merch"})
			return
		}

		var body struct {
			Name     *string  `json:"name"`
			Price    *float64 `json:"price"`
			Discount *float64 `json:"discount"`
			Stock    *int     `json:"stock"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.RespondWithJSON(w, 400, map[string]any{"success": false, "error": "invalid json"})
			return
		}

		update := map[string]any{
			"updatedat": time.Now(),
		}

		if body.Name != nil {
			update["name"] = *body.Name
		}
		if body.Price != nil && *body.Price > 0 {
			update["price"] = *body.Price
		}
		if body.Discount != nil {
			update["discount"] = *body.Discount
		}
		if body.Stock != nil && *body.Stock >= 0 {
			update["stock"] = *body.Stock
		}

		if err := updateMerchFields(r.Context(), app, entityType, eventID, merchID, update); err != nil {
			utils.RespondWithJSON(w, 404, map[string]any{"success": false, "error": "merch not found"})
			return
		}

		// Log audit trail for merchandise update
		auditlog.LogAction(
			r.Context(), app, r, userID,
			auditlog.AuditActionMerchUpdate,
			"merchandise", merchID, "success",
			map[string]interface{}{
				"updates": update,
			},
		)

		mqpayload, _ := json.Marshal(mqevent.MerchUpdatedPayload{})

		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.MerchUpdatedEvent, mqpayload)

		utils.RespondWithJSON(w, 200, map[string]any{"success": true})
	}
}

// ---------------------- Delete Merch ----------------------
func DeleteMerch(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		entityType := utils.GetParam(r, "entityType")
		eventID := utils.GetParam(r, "eventid")
		merchID := utils.GetParam(r, "merchid")

		// SECURITY: Verify user is authenticated
		userID, ok := r.Context().Value(config.UserIDKey).(string)
		if !ok || userID == "" {
			utils.RespondWithJSON(w, 401, map[string]any{"success": false, "error": "unauthorized"})
			return
		}

		// SECURITY: Verify user is the owner of the entity
		if owner, err := getEntityOwner(r.Context(), app, entityType, eventID); err != nil {
			if err.Error() == "cannot verify ownership" {
				utils.RespondWithJSON(w, 403, map[string]any{"success": false, "error": "cannot verify ownership"})
				return
			}
			utils.RespondWithJSON(w, 404, map[string]any{"success": false, "error": "entity not found"})
			return
		} else if owner != userID {
			utils.RespondWithJSON(w, 403, map[string]any{"success": false, "error": "forbidden: only entity owner can delete merch"})
			return
		}

		// SECURITY: Use soft delete instead of hard delete
		now := time.Now()
		if err := softDeleteMerch(r.Context(), app, entityType, eventID, merchID, now); err != nil {
			utils.RespondWithJSON(w, 404, map[string]any{"success": false, "error": "merch not found"})
			return
		}

		// Log audit trail for merchandise deletion
		auditlog.LogAction(
			r.Context(), app, r, userID,
			auditlog.AuditActionMerchDelete,
			"merchandise", merchID, "success",
			map[string]interface{}{
				"deleted_at": now,
			},
		)

		mqpayload, _ := json.Marshal(mqevent.MerchDeletedPayload{})

		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.MerchDeletedEvent, mqpayload)

		utils.RespondWithJSON(w, 200, map[string]any{"success": true})
	}
}

// ---------------------- Buy Merch ----------------------
func BuyMerch(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := r.Context().Value(config.UserIDKey).(string)
		if !ok || userID == "" {
			utils.RespondWithJSON(w, 401, map[string]any{"success": false, "error": "unauthorized"})
			return
		}

		var body struct {
			Quantity int `json:"quantity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Quantity < 1 {
			utils.RespondWithJSON(w, 400, map[string]any{"success": false, "error": "invalid quantity"})
			return
		}

		err := buyMerchTransaction(r.Context(), app, r.Context(), userID, utils.GetParam(r, "entityType"), utils.GetParam(r, "eventid"), utils.GetParam(r, "merchid"), body.Quantity)

		if err != nil {
			utils.RespondWithJSON(w, 400, map[string]any{"success": false, "error": err.Error()})
			return
		}

		mqpayload, _ := json.Marshal(mqevent.MerchBoughtPayload{})

		_ = mq.PublishWithMeta(ctx, app.MQ, mqevent.MerchBoughtEvent, mqpayload)

		utils.RespondWithJSON(w, 200, map[string]any{
			"success": true,
			"message": "purchase successful",
		})
	}
}
