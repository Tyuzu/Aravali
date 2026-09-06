package deliveries

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"scav/infra"
	"scav/utils"
)

// Handler for generic status updates (e.g. PATCH /api/v1/deliveries/:deliveryid/status)
func UpdateDeliveryStatus(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")

		var req struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		delivery, err := updateDeliveryStatus(app, r, deliveryID, req.Status)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, delivery)
	}
}

func CancelDelivery(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")
		delivery, err := updateDeliveryStatus(app, r, deliveryID, StatusCancelled)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, delivery)
	}
}

func AssignDriver(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")
		tenantID := GetTenantIDFromContext(r.Context())

		var req struct {
			DriverID string `json:"driver_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		ctx := r.Context()
		updated, err := upsertDeliveryAssignment(ctx, app, deliveryID, tenantID, req.DriverID, utils.GetUserIDFromRequest(r))
		if err != nil {
			if err.Error() == "delivery not found" {
				utils.RespondWithError(w, http.StatusNotFound, "Delivery not found")
				return
			}
			utils.RespondWithError(w, http.StatusConflict, err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, updated)
	}
}

func AcceptAssignment(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")
		delivery, err := updateDeliveryStatus(app, r, deliveryID, StatusAccepted)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, delivery)
	}
}

func MarkPickedUp(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")
		delivery, err := updateDeliveryStatus(app, r, deliveryID, StatusPickedUp)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, delivery)
	}
}

func StartDelivery(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")
		delivery, err := updateDeliveryStatus(app, r, deliveryID, StatusInTransit)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, delivery)
	}
}

func CompleteDelivery(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")
		delivery, err := updateDeliveryStatus(app, r, deliveryID, StatusDelivered)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, delivery)
	}
}

func updateDeliveryStatus(app *infra.Deps, r *http.Request, deliveryID string, newStatus string) (*Delivery, error) {
	ctx := r.Context()
	userID := utils.GetUserIDFromRequest(r)
	tenantID := GetTenantIDFromContext(ctx)
	return updateDeliveryStatusRecord(ctx, app, deliveryID, tenantID, userID, newStatus)
}

func CreateDelivery(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			PickupLoc  Location `json:"pickup_loc"`
			DropoffLoc Location `json:"dropoff_loc"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		ctx := r.Context()
		now := time.Now()
		deliveryID := utils.GenerateRandomString(18)

		delivery := Delivery{
			DeliveryID:          deliveryID,
			UserID:              utils.GetUserIDFromRequest(r),
			TenantID:            GetTenantIDFromContext(ctx),
			Status:              StatusCreated,
			PickupLoc:           req.PickupLoc,
			DropoffLoc:          req.DropoffLoc,
			PublicTrackingToken: utils.GenerateRandomString(32),
			CreatedAt:           now,
			UpdatedAt:           now,
			StatusHistory: []StatusHistoryItem{
				{
					Status:    StatusCreated,
					Timestamp: now,
					UpdatedBy: utils.GetUserIDFromRequest(r),
				},
			},
		}

		if err := saveDelivery(ctx, app, delivery); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create delivery")
			return
		}

		_ = app.NatsConn.Publish("deliveries.created", []byte(delivery.DeliveryID))
		utils.RespondWithJSON(w, http.StatusCreated, delivery)
	}
}

func GetMyDeliveries(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := utils.GetUserIDFromRequest(r)
		tenantID := GetTenantIDFromContext(ctx)

		myDeliveries, err := findMyDeliveries(ctx, app, userID, tenantID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch deliveries")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, myDeliveries)
	}
}

func GetDeliveryByID(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")
		tenantID := GetTenantIDFromContext(r.Context())
		if deliveryID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Missing delivery ID")
			return
		}

		ctx := r.Context()
		cacheKey := fmt.Sprintf("delivery:%s", deliveryID)

		if cached, err := app.Cache.Get(ctx, cacheKey); err == nil && len(cached) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(cached)
			return
		}

		delivery, err := fetchDeliveryForRead(ctx, app, deliveryID, tenantID)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Delivery not found")
			return
		}

		if bytes, err := json.Marshal(delivery); err == nil {
			_ = app.Cache.Set(ctx, cacheKey, bytes, 15*time.Minute)
		}

		utils.RespondWithJSON(w, http.StatusOK, delivery)
	}
}
