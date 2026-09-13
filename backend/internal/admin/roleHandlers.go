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

func ApplyForRole(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		var payload ApplyPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		application, err := ProcessApplyForRole(r.Context(), app, userID, payload)
		if err != nil {
			switch {
			case errors.Is(err, ErrInvalidPayload):
				utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			case errors.Is(err, ErrUserRoleDefault):
				utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			case errors.Is(err, ErrPendingRequestExists):
				utils.RespondWithError(w, http.StatusConflict, err.Error())
			default:
				utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save role request")
			}
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, map[string]any{
			"message": "Role request submitted",
			"id":      application.ID,
			"status":  application.Status,
		})
	}
}

func GetMyRoleRequests(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		applications, err := FetchUserRoleRequests(r.Context(), app, userID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to load your role requests")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, applications)
	}
}

func ListRoleRequests(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawStatus := r.URL.Query().Get("status")
		applications, err := FetchAllRoleRequests(r.Context(), app, rawStatus)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to load role applications")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, applications)
	}
}

func ApproveRoleRequest(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := utils.GetParam(r, "id")
		if appID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Missing application id")
			return
		}

		approvedRole, err := ProcessApproveRoleRequest(r.Context(), app, appID)
		if err != nil {
			switch {
			case errors.Is(err, ErrApplicationNotFound):
				utils.RespondWithError(w, http.StatusNotFound, "Application not found")
			case errors.Is(err, ErrAlreadyResolved):
				utils.RespondWithError(w, http.StatusConflict, "This role request has already been resolved")
			case errors.Is(err, ErrUserNotFound):
				utils.RespondWithError(w, http.StatusNotFound, "User not found")
			default:
				utils.RespondWithError(w, http.StatusInternalServerError, "Failed to approve role request")
			}
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"message": "Role approved",
			"role":    approvedRole,
		})
	}
}

func RejectRoleRequest(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := utils.GetParam(r, "id")
		if appID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Missing application id")
			return
		}

		err := ProcessRejectRoleRequest(r.Context(), app, appID)
		if err != nil {
			switch {
			case errors.Is(err, ErrApplicationNotFound):
				utils.RespondWithError(w, http.StatusNotFound, "Application not found")
			case errors.Is(err, ErrAlreadyResolved):
				utils.RespondWithError(w, http.StatusConflict, "This role request has already been resolved")
			default:
				utils.RespondWithError(w, http.StatusInternalServerError, "Failed to reject role request")
			}
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Role request rejected",
		})
	}
}
