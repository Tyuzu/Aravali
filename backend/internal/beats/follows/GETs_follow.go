package follows

import (
	"net/http"

	"scav/infra"
	"scav/internal/auth"
	"scav/utils"
)

// GET /api/v1/follow/:id
func DoesFollow(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := utils.GetUserIDFromRequest(r)
		followedUserID := utils.GetParam(r, "id")

		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if followedUserID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		count, err := CountFollowRelationship(r.Context(), app, userID, followedUserID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := map[string]bool{
			"isFollowing": count > 0,
		}

		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}

// GET /api/v1/followers
func GetFollowers(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var userFollow UserFollow
		err := FindFollowEntryByUserID(r.Context(), app, userID, &userFollow)
		if err != nil || len(userFollow.Followers) == 0 {
			utils.RespondWithJSON(w, http.StatusOK, []auth.User{})
			return
		}

		var followers []auth.User
		err = FindUsersByIDsForFollow(r.Context(), app, userFollow.Followers, &followers)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, followers)
	}
}

// GET /api/v1/following
func GetFollowing(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var userFollow UserFollow
		err := FindFollowEntryByUserID(r.Context(), app, userID, &userFollow)
		if err != nil || len(userFollow.Follows) == 0 {
			utils.RespondWithJSON(w, http.StatusOK, []auth.User{})
			return
		}

		var following []auth.User
		err = FindUsersByIDsForFollow(r.Context(), app, userFollow.Follows, &following)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, following)
	}
}
