package profile

import (
	"net/http"
	"strings"

	"scav/infra"
	"scav/utils"

	"github.com/julienschmidt/httprouter"
)

func GetProfile(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		username := utils.GetUsernameFromRequest(r)
		if username == "" {
			username = "user_" + userID
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"userid":   userID,
			"username": username,
			"email":    "",
			"name":     "",
			"bio":      "",
		})
	}
}

func EditProfile(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var payload ProfileUpdateRequest
		if err := utils.ParseJSON(r, &payload); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		username := strings.TrimSpace(payload.Username)
		if username == "" {
			username = utils.GetUsernameFromRequest(r)
		}
		if username == "" {
			username = "user_" + userID
		}

		updated := map[string]any{
			"userid":       userID,
			"username":     username,
			"email":        strings.TrimSpace(payload.Email),
			"name":         strings.TrimSpace(payload.Name),
			"bio":          strings.TrimSpace(payload.Bio),
			"phone_number": strings.TrimSpace(payload.PhoneNumber),
			"updated_at":   r.URL.Query().Get("updated_at"),
		}
		utils.RespondWithJSON(w, http.StatusOK, updated)
	}
}

func DeleteProfile(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if utils.GetUserIDFromRequest(r) == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetUserProfile(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		username := strings.TrimSpace(params.ByName("username"))
		if username == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "username is required")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"userid":    "public_" + username,
			"username":  username,
			"email":     "",
			"bio":       "",
			"name":      "",
			"is_public": true,
		})
	}
}
