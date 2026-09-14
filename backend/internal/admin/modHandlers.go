package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"scav/infra"
	"scav/utils"
)

// ============================================================================
// HTTP Handlers
// ============================================================================

func ApplyModerator(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload ApplyModeratorPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, `{"error":"Invalid JSON payload"}`, http.StatusBadRequest)
			return
		}

		appx, err := ProcessApplyModerator(r.Context(), app, payload)
		if err != nil {
			switch {
			case errors.Is(err, ErrMissingRequiredFields):
				http.Error(w, `{"error":"Missing required fields"}`, http.StatusBadRequest)
			case errors.Is(err, ErrAlreadyApplied):
				http.Error(w, `{"error":"You have already applied to be a moderator"}`, http.StatusConflict)
			default:
				http.Error(w, `{"error":"Failed to save application"}`, http.StatusInternalServerError)
			}
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message": "Moderator application submitted",
			"id":      appx.ID,
		})
	}
}

func ListModeratorApplications(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		applications, err := FetchModeratorApplications(r.Context(), app, status)
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
		id := utils.GetParam(r, "id")
		if id == "" {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid or missing application ID",
			})
			return
		}

		if err := ProcessApproveModerator(r.Context(), app, id); err != nil {
			if errors.Is(err, ErrModAppNotFound) {
				utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{
					"error": "Application not found or update failed",
				})
				return
			}
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Failed to approve application",
			})
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Application approved successfully",
		})
	}
}

func RejectModerator(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := utils.GetParam(r, "id")
		if id == "" {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid or missing application ID",
			})
			return
		}

		if err := ProcessRejectModerator(r.Context(), app, id); err != nil {
			if errors.Is(err, ErrModAppNotFound) {
				utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{
					"error": "Application not found or update failed",
				})
				return
			}
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Failed to reject application",
			})
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Application rejected successfully",
		})
	}
}
