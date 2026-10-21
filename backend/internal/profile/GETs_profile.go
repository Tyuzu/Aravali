// File: internal/profile/GETs_profile.go

package profile

import (
	"encoding/json"
	"net/http"
	"time"

	"scav/infra"
	"scav/utils"
)

/* -------------------------------------------------------
   Get Own Profile
------------------------------------------------------- */

func GetProfile(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestingUserID := utils.GetUserIDFromRequest(r)

		user, err := FindUserByFilter(ctx, app, "userid = $1", []any{requestingUserID})
		if err != nil || user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		user.Online, _ = isOnline(ctx, user.UserID, app.Cache)

		profileJSON, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Encoding failed", http.StatusInternalServerError)
			return
		}

		// Best-effort cache write (5 min TTL)
		_ = CacheProfile(ctx, app.Cache, user.Username, string(profileJSON), 5*time.Minute)

		utils.RespondWithJSON(w, http.StatusOK, user)
	}
}

/* -------------------------------------------------------
   Get Another User's Profile
------------------------------------------------------- */

func GetUserProfile(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		username := utils.GetParam(r, "username")

		user, err := FindUserByFilter(ctx, app, "username = $1", []any{username})
		if err != nil || user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		online, _ := isOnline(ctx, user.UserID, app.Cache)

		response := UserProfileResponse{
			UserID:    user.UserID,
			Username:  user.Username,
			Email:     user.Email,
			Name:      user.Name,
			Bio:       user.Bio,
			Avatar:    user.Avatar,
			Banner:    user.Banner,
			Online:    online,
			LastLogin: user.LastLogin,
		}

		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}
