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

		farms, total, err := getMyFarmsPage(ctx, app.DB, userID, skip, limit)
		if err != nil {
			utils.RespondWithError(
				w,
				http.StatusInternalServerError,
				"Error fetching farms",
			)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"farms":   farms,
			"total":   total,
			"page":    skip/limit + 1,
			"limit":   limit,
		})
	}
}
