// File: internal/beats/suggestions/suggestions.go

package suggestions

import (
	"context"
	"net/http"
	"scav/internal/places"
	"strconv"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	log "scav/infra/logger"
	"scav/utils"
)

func SuggestFollowers(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserID := r.URL.Query().Get("userid")
		if currentUserID == "" {
			http.Error(w, "Missing userid", http.StatusBadRequest)
			return
		}

		userID, ok := r.Context().Value(config.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit < 1 {
			limit = 10
		}

		offset := int64((page - 1) * limit)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		followData, err := findFollowDataByUserID(ctx, app, currentUserID)
		if err != nil {
			followData.Follows = []string{}
		}

		excludedUserIDs := append(followData.Follows, currentUserID, userID)

		where := "1 = 1"
		args := []any{}
		if len(excludedUserIDs) > 0 {
			placeholders := make([]string, len(excludedUserIDs))
			args = make([]any, len(excludedUserIDs))
			for i := range excludedUserIDs {
				placeholders[i] = "$" + strconv.Itoa(i+1)
				args[i] = excludedUserIDs[i]
			}
			where = "userid NOT IN (" + strings.Join(placeholders, ", ") + ")"
		}

		users, err := findSuggestedUsers(ctx, app, where, args)
		if err != nil {
			http.Error(w, "Failed to fetch suggestions", http.StatusInternalServerError)
			return
		}

		for i := range users {
			users[i].IsFollowing = false
		}

		utils.SortAndSlice(
			&users,
			[]utils.SortField{{Key: "userid", Value: 1}},
			offset,
			int64(limit),
		)

		if users == nil {
			users = []UserSuggest{}
		}

		utils.RespondWithJSON(w, http.StatusOK, users)
	}
}

func GetNearbyPlaces(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		curplace := r.URL.Query().Get("place")
		if len(curplace) != 14 {
			log.Printf("invalid place id: %s", curplace)
		}

		nearbyplaces, err := findNearbyPlaces(ctx, app, "1 = 1", nil)
		if err != nil {
			http.Error(w, "Failed to fetch places", http.StatusInternalServerError)
			return
		}

		if nearbyplaces == nil {
			nearbyplaces = []places.Place{}
		}

		sanitized := make([]map[string]any, 0, len(nearbyplaces))
		for _, place := range nearbyplaces {
			if place.PlaceID == curplace {
				continue
			}
			sanitized = append(sanitized, map[string]any{
				"placeid":     place.PlaceID,
				"name":        place.Name,
				"banner":      place.Banner,
				"category":    place.Category,
				"capacity":    place.Capacity,
				"reviewCount": place.ReviewCount,
			})
		}

		utils.RespondWithJSON(w, http.StatusOK, sanitized)
	}
}
