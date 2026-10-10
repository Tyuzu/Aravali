// File: internal/workers/GETs_workers.go

package workers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"scav/infra"
	log "scav/infra/logger"
	"scav/utils"
)

/* -------------------- Workers -------------------- */

func GetWorkerById(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		workerID := utils.GetParam(r, "workerId")
		worker, err := findWorkerByIDFromDB(ctx, app, workerID)
		if err != nil {
			log.Printf("DB error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch worker")
			return
		}

		if worker.BaitoWorkerId == "" {
			utils.RespondWithError(w, http.StatusNotFound, "Worker not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, worker)
	}
}

func GetWorkerSkills(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		skills, err := getUniqueWorkerSkillsFromDB(ctx, app)
		if err != nil {
			log.Printf("Aggregate error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch skills")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, skills)
	}
}

// GetWorkers returns a list of workers with optional search and skill filtering.
func GetWorkers(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		search := strings.TrimSpace(r.URL.Query().Get("search"))
		skill := strings.TrimSpace(r.URL.Query().Get("skill"))
		skip, limit := utils.ParsePagination(r, 10, 100)

		workers, total, err := getWorkersPage(ctx, app, search, skill, skip, limit)
		if err != nil {
			log.Printf("DB error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch workers")
			return
		}

		if workers == nil {
			workers = []BaitoWorkersResponse{}
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"data":  workers,
			"total": total,
		})
	}
}
