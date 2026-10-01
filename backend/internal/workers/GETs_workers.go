// File: internal/workers/GETs_workers.go

package workers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"scav/infra"
	"scav/infra/sqldb"
	"scav/utils"
	log "scav/utils/logger"
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

func buildWorkerQuery(search string, skill string) (string, []any) {
	clauses := make([]string, 0, 2)
	args := make([]any, 0, 4)

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		namePH := fmt.Sprintf("$%d", len(args)+1)
		args = append(args, searchTerm)
		locPH := fmt.Sprintf("$%d", len(args)+1)
		args = append(args, searchTerm)
		bioPH := fmt.Sprintf("$%d", len(args)+1)
		args = append(args, searchTerm)
		clauses = append(clauses, fmt.Sprintf("(LOWER(name) LIKE LOWER(%s) OR LOWER(location) LIKE LOWER(%s) OR LOWER(bio) LIKE LOWER(%s))", namePH, locPH, bioPH))
	}

	if skill != "" {
		skillTerm := "%" + strings.ToLower(skill) + "%"
		skillPH := fmt.Sprintf("$%d", len(args)+1)
		args = append(args, skillTerm)
		clauses = append(clauses, fmt.Sprintf("(LOWER(CAST(preferred AS TEXT)) LIKE LOWER(%s))", skillPH))
	}

	if len(clauses) == 0 {
		return "TRUE", nil
	}

	return strings.Join(clauses, " AND "), args
}

// GetWorkers returns a list of workers with optional search and skill filtering.
func GetWorkers(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		search := strings.TrimSpace(r.URL.Query().Get("search"))
		skill := strings.TrimSpace(r.URL.Query().Get("skill"))

		whereClause, args := buildWorkerQuery(search, skill)
		skip, limit := utils.ParsePagination(r, 10, 100)

		opts := sqldb.FindManyOptions{
			Limit:  int64(limit),
			Offset: int64(skip),
			Sort: []sqldb.OrderBy{{
				Column:     "createdat",
				Descending: true,
			}},
		}

		workers, err := findWorkersFromDB(ctx, app, whereClause, args, opts)
		if err != nil {
			log.Printf("DB error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch workers")
			return
		}

		total, err := countWorkersFromDB(ctx, app, whereClause, args)
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
