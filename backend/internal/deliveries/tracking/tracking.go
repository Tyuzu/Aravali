package tracking

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"scav/infra"
	"scav/internal/deliveries"
	"scav/utils"
)

func GetDeliveryTracking(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		ctx := r.Context()

		result, err := getTrackingDetails(ctx, app, deliveryID, tenantID)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Tracking details not found")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, result)
	}
}

func GetDeliveryLocation(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")
		ctx := r.Context()

		locationBytes, err := app.Cache.Get(ctx, fmt.Sprintf("location:delivery:%s", deliveryID))
		if err == nil && len(locationBytes) > 0 {
			var loc deliveries.GPSData
			_ = json.Unmarshal(locationBytes, &loc)
			utils.RespondWithJSON(w, http.StatusOK, loc)
			return
		}

		utils.RespondWithError(w, http.StatusNotFound, "Location not available")
	}
}

func GetDeliveryEvents(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		ctx := r.Context()

		events, err := getDeliveryEvents(ctx, app, deliveryID, tenantID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve events")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, events)
	}
}

func GetStatusHistory(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		ctx := r.Context()

		history, err := getStatusHistory(ctx, app, deliveryID, tenantID)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "History not found")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, history)
	}
}

func AddProof(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")
		tenantID := deliveries.GetTenantIDFromContext(r.Context())

		var req struct {
			ProofType string `json:"type"`
			URL       string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid body")
			return
		}

		ctx := r.Context()
		proof := deliveries.Proof{
			ProofID:   utils.GenerateRandomString(16),
			Type:      req.ProofType,
			URL:       req.URL,
			CreatedAt: time.Now(),
		}

		if err := addProofToDelivery(ctx, app, deliveryID, tenantID, proof); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add proof")
			return
		}
		utils.RespondWithJSON(w, http.StatusCreated, proof)
	}
}

func GetProof(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := utils.GetParam(r, "deliveryid")
		tenantID := deliveries.GetTenantIDFromContext(r.Context())
		ctx := r.Context()

		proofs, err := getProofs(ctx, app, deliveryID, tenantID)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Proof not found")
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, proofs)
	}
}

func GetPublicTracking(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetParam(r, "token")
		ctx := r.Context()

		res, err := getPublicTrackingInfo(ctx, app, token)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Invalid or expired tracking token")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, res)
	}
}

func GetPublicLocation(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.GetParam(r, "token")
		ctx := r.Context()

		deliveryIDBytes, err := app.Cache.Get(ctx, fmt.Sprintf("token:map:%s", token))
		if err != nil || len(deliveryIDBytes) == 0 {
			utils.RespondWithError(w, http.StatusNotFound, "Tracking token location unavailable")
			return
		}

		locBytes, err := app.Cache.Get(ctx, fmt.Sprintf("location:delivery:%s", string(deliveryIDBytes)))
		if err != nil || len(locBytes) == 0 {
			utils.RespondWithError(w, http.StatusNotFound, "Live location not streaming")
			return
		}

		var loc deliveries.GPSData
		_ = json.Unmarshal(locBytes, &loc)
		utils.RespondWithJSON(w, http.StatusOK, loc)
	}
}
