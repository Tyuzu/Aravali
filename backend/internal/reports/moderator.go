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

		var reports []Report
		err := FindReports(ctx, app, filter, &reports)
		if err != nil {
			http.Error(w, `{"error":"Failed to fetch reports"}`, http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, reports)
	}
}
