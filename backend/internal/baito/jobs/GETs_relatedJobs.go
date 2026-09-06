package jobs

import (
	"net/http"
	"scav/infra"
	"scav/internal/baito"
	"scav/utils"
)

// ------------------ READ LIST ------------------
func GetJobsRelatedTOEntity(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		entityType := utils.GetParam(r, "entitytype")
		entityID := utils.GetParam(r, "entityid")

		jobs, err := FindJobsForEntity(ctx, app, entityType, entityID)
		if err != nil {
			http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
			return
		}

		if jobs == nil {
			jobs = []baito.BaitosResponse{}
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{"jobs": jobs})
	}
}
