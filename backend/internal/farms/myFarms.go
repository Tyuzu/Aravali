package farms

import (
	"context"
	"net/http"
	"scav/infra"
	"scav/utils"
	"time"
)

func GetMyFarms(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		userID := utils.GetUserIDFromRequest(r)

		skip, limit := utils.ParsePagination(r, 10, 100)

		pipeline := []any{
			map[string]any{
				"$match": map[string]any{
					"createdBy": userID,
				},
			},
			map[string]any{
				"$sort": map[string]any{
					"createdAt": -1,
				},
			},
			map[string]any{
				"$lookup": map[string]any{
					"from":         "crops",
					"localField":   "farmid",
					"foreignField": "farmid",
					"as":           "crops",
				},
			},
			map[string]any{"$skip": skip},
			map[string]any{"$limit": limit},
		}

		var farms []Farm

		if err := app.DB.Aggregate(
			ctx,
			farmsCollection,
			pipeline,
			&farms,
		); err != nil {
			utils.RespondWithError(
				w,
				http.StatusInternalServerError,
				"Error fetching farms",
			)
			return
		}

		total, _ := app.DB.CountDocuments(
			ctx,
			farmsCollection,
			map[string]any{
				"createdBy": userID,
			},
		)

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"farms":   farms,
			"total":   total,
			"page":    skip/limit + 1,
			"limit":   limit,
		})
	}
}
