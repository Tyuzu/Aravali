// File: internal/reports/moderator.go

package reports

import (
	"net/http"

	"scav/infra"
	"scav/utils"
)

func GetReportsForMod(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		filter := map[string]any{
			"status": map[string]any{
				"$nin": []string{"resolved", "rejected"},
			},
		}

		query, args := buildFilterQuery(filter)
		var reports []Report
		err := FindReports(ctx, app, query, args, &reports)
		if err != nil {
			http.Error(w, `{"error":"Failed to fetch reports"}`, http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, reports)
	}
}
