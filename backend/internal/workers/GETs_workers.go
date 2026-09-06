package workers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"scav/infra"
	"scav/infra/db"
	"scav/utils"
	log "scav/utils/logger"

	"go.mongodb.org/mongo-driver/bson"
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

		filter := map[string]any{}

		if search != "" {
			filter["$or"] = []any{
				map[string]any{"name_contains": search},
				map[string]any{"location_contains": search},
				map[string]any{"bio_contains": search},
			}
		}

		if skill != "" {
			filter["preferredRoles"] = skill
		}

		skip, limit := utils.ParsePagination(r, 10, 100)

		opts := db.FindManyOptions{
			Skip:  skip,
			Limit: limit,
			Sort:  []bson.E{{Key: "createdAt", Value: -1}},
		}

		workers, err := findWorkersFromDB(ctx, app, filter, opts)
		if err != nil {
			log.Printf("DB error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch workers")
			return
		}

		total, err := countWorkersFromDB(ctx, app, filter)
		if err != nil {
			log.Printf("Count error: %v", err)
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
