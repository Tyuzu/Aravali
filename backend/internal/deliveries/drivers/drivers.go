package drivers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"scav/infra"
	"scav/internal/deliveries"
	"scav/utils"
)

// Helper to resolve driver ID from context or request fallback
func resolveDriverID(r *http.Request) string {
	if driverID := deliveries.GetDriverIDFromContext(r.Context()); driverID != "" {
		return driverID
	}
	return utils.GetUserIDFromRequest(r)
}

func GetProfile(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		driverID := resolveDriverID(r)
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		ctx := r.Context()

		driver, err := getDriverProfileByID(ctx, app, driverID, tenantID)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Driver profile not found")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, driver)
	}
}

func UpdateProfile(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		driverID := resolveDriverID(r)
		tenantID := deliveries.GetTenantIDFromContext(r.Context())

		var updates map[string]any
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid body")
			return
		}

		ctx := r.Context()
		delete(updates, "id")
		delete(updates, "tenantid")
		updates["updated_at"] = time.Now()

		if err := updateDriverProfile(ctx, app, driverID, tenantID, updates); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update driver")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

func GoOnline(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		driverID := resolveDriverID(r)
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		ctx := r.Context()

		_ = app.Cache.HSet(ctx, "drivers:online", driverID, []byte("true"))
		_ = setDriverOnlineState(ctx, app, driverID, tenantID, true)

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "online"})
	}
}

func GoOffline(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		driverID := resolveDriverID(r)
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		ctx := r.Context()

		_, _ = app.Cache.HDel(ctx, "drivers:online", driverID)
		_ = setDriverOnlineState(ctx, app, driverID, tenantID, false)

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "offline"})
	}
}

func GetStatus(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		driverID := resolveDriverID(r)
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		ctx := r.Context()

		status, err := getDriverStatus(ctx, app, driverID, tenantID)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Driver status unavailable")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, status)
	}
}

func GetAvailableJobs(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		ctx := r.Context()

		jobs, err := getAvailableJobsForTenant(ctx, app, tenantID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch available jobs")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, jobs)
	}
}

func GetActiveDeliveries(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		driverID := resolveDriverID(r)
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		ctx := r.Context()

		active, err := getActiveJobsForDriver(ctx, app, driverID, tenantID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch active deliveries")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, active)
	}
}

// ClaimJob allows a driver to pick up an unassigned job (POST /drivers/me/deliveries/:deliveryid/claim)
func ClaimJob(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		driverID := resolveDriverID(r)
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		deliveryID := utils.GetParam(r, "deliveryid")
		ctx := r.Context()

		updated, err := claimDeliveryAssignment(ctx, app, deliveryID, tenantID, driverID)
		if err != nil {
			if err.Error() == "delivery is already assigned to another driver" {
				utils.RespondWithError(w, http.StatusConflict, err.Error())
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to claim delivery")
			return
		}

		_ = app.Cache.Del(ctx, fmt.Sprintf("delivery:%s", deliveryID))
		utils.RespondWithJSON(w, http.StatusOK, updated)
	}
}

func AcceptJob(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		driverID := resolveDriverID(r)
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		deliveryID := utils.GetParam(r, "deliveryid")
		ctx := r.Context()

		delivery, err := acceptDeliveryAssignment(ctx, app, deliveryID, tenantID, driverID)
		if err != nil {
			utils.RespondWithError(w, http.StatusConflict, "Job no longer available or invalid")
			return
		}

		_ = app.Cache.Del(ctx, fmt.Sprintf("delivery:%s", deliveryID))
		_ = app.NatsConn.Publish(fmt.Sprintf("deliveries.status.%s", deliveries.StatusAccepted), []byte(deliveryID))

		utils.RespondWithJSON(w, http.StatusOK, delivery)
	}
}

func RejectJob(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		driverID := resolveDriverID(r)
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		deliveryID := utils.GetParam(r, "deliveryid")
		ctx := r.Context()

		_ = saveDriverRejection(ctx, app, tenantID, driverID, deliveryID)

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "rejected"})
	}
}

func SendGPSLocation(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		driverID := resolveDriverID(r)
		var loc deliveries.GPSData
		if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid payload")
			return
		}
		loc.Timestamp = time.Now()

		ctx := r.Context()
		bytes, _ := json.Marshal(loc)

		_ = app.Cache.Set(ctx, fmt.Sprintf("gps:driver:%s", driverID), bytes, 1*time.Hour)
		_ = app.NatsConn.Publish(fmt.Sprintf("drivers.location.%s", driverID), bytes)

		utils.RespondWithJSON(w, http.StatusAccepted, map[string]string{"status": "location_updated"})
	}
}

func GetCurrentGPS(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		driverID := resolveDriverID(r)
		ctx := r.Context()

		val, err := app.Cache.Get(ctx, fmt.Sprintf("gps:driver:%s", driverID))
		if err != nil || len(val) == 0 {
			utils.RespondWithError(w, http.StatusNotFound, "No recent GPS data found")
			return
		}

		var loc deliveries.GPSData
		_ = json.Unmarshal(val, &loc)
		utils.RespondWithJSON(w, http.StatusOK, loc)
	}
}
