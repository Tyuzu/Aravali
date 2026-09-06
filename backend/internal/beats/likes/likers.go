package likes

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"scav/infra"
	"scav/utils"
)

func GetLikers(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(
			r.Context(),
			3*time.Second,
		)
		defer cancel()

		entityType := utils.GetParam(r, "entitytype")
		entityID := utils.GetParam(r, "entityid")

		if entityType == "" || entityID == "" {
			http.Error(
				w,
				"Bad Request",
				http.StatusBadRequest,
			)
			return
		}

		if !IsValidEntityType(entityType) {
			http.Error(
				w,
				"Invalid entity type",
				http.StatusBadRequest,
			)
			return
		}

		limit := int64(20)

		if value := r.URL.Query().Get("limit"); value != "" {
			parsed, err := strconv.ParseInt(value, 10, 64)
			if err != nil || parsed <= 0 {
				http.Error(
					w,
					"Invalid limit",
					http.StatusBadRequest,
				)
				return
			}

			if parsed > 100 {
				parsed = 100
			}

			limit = parsed
		}

		likes, err := FindLikesByEntity(ctx, app, entityType, entityID)
		if err != nil {
			http.Error(
				w,
				"Failed to fetch likers",
				http.StatusInternalServerError,
			)
			return
		}

		if int64(len(likes)) > limit {
			likes = likes[:limit]
		}

		if len(likes) == 0 {
			utils.RespondWithJSON(
				w,
				http.StatusOK,
				map[string]any{
					"likers": []any{},
				},
			)
			return
		}

		userIDs := make([]string, 0, len(likes))

		for _, like := range likes {
			userIDs = append(userIDs, like.UserID)
		}

		users, err := FindUsersByIDs(ctx, app, userIDs)
		if err != nil {
			http.Error(
				w,
				"Failed to fetch users",
				http.StatusInternalServerError,
			)
			return
		}

		userMap := make(map[string]map[string]string, len(users))

		for _, user := range users {
			userMap[user["userid"]] = user
		}

		likers := make([]map[string]string, 0, len(likes))

		for _, like := range likes {
			if user, ok := userMap[like.UserID]; ok {
				likers = append(likers, user)
			}
		}

		utils.RespondWithJSON(
			w,
			http.StatusOK,
			map[string]any{
				"likers": likers,
			},
		)
	}
}
