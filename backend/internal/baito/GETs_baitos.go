// File: internal/baito/GETs_baitos.go

package baito

import (
	"context"
	"net/http"

	"scav/infra"
	"scav/infra/logger"
	log "scav/infra/logger"
	"scav/utils"
)

/* -------------------- Helpers -------------------- */

func enrichBaitoApplicationCount(ctx *context.Context, app *infra.Deps, baito *Baito) error {
	count, err := countApplicationsForBaito(*ctx, app, baito.BaitoId)
	if err != nil {
		return err
	}

	baito.ApplicationCount = int(count)
	return nil
}

func respondBaitos(w http.ResponseWriter, baitos []BaitosResponse) {
	if baitos == nil {
		baitos = []BaitosResponse{}
	}
	utils.RespondWithJSON(w, http.StatusOK, baitos)
}

/* -------------------- Handlers -------------------- */

func GetLatestBaitos(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		baitos, err := findLatestBaitosFromDB(ctx, app, "1 = 1", nil, 20)
		if err != nil {
			logger.Printf("DB error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
			return
		}

		respondBaitos(w, baitos)
	}
}

func GetRelatedBaitos(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		category := r.URL.Query().Get("category")
		exclude := r.URL.Query().Get("exclude")

		query := "1 = 1"
		args := []any{}
		if category != "" {
			query = "category = $1"
			args = append(args, category)
		}
		if exclude != "" {
			if query == "1 = 1" {
				query = "baitoid <> $1"
			} else {
				query += " AND baitoid <> $2"
			}
			args = append(args, exclude)
		}

		baitos, err := findRelatedBaitosFromDB(ctx, app, query, args, 10)
		if err != nil {
			logger.Printf("DB error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
			return
		}

		// fallback if none found
		if len(baitos) == 0 {
			fallbackQuery := "1 = 1"
			fallbackArgs := []any{}
			if exclude != "" {
				fallbackQuery = "baitoid <> $1"
				fallbackArgs = append(fallbackArgs, exclude)
			}

			baitos, _ = findRelatedBaitosFromDB(ctx, app, fallbackQuery, fallbackArgs, 10)
		}

		respondBaitos(w, baitos)
	}
}
func GetBaitoByID(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := utils.GetParam(r, "baitoid")

		b, err := findBaitoByIDFromDB(ctx, app, id)
		if err != nil {
			logger.Printf("DB error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
			return
		}

		// Check if zero-value struct was returned by repository layer
		if b.BaitoId == "" {
			utils.RespondWithError(w, http.StatusNotFound, "Not found")
			return
		}

		if err := enrichBaitoApplicationCount(&ctx, app, &b); err != nil {
			logger.Printf("Failed to count applications for baito %s: %v", id, err)
		}

		utils.RespondWithJSON(w, http.StatusOK, b)
	}
}

func GetMyBaitos(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := utils.GetUserIDFromRequest(r)

		baitos, err := findMyBaitosFromDB(ctx, app, userID)
		if err != nil {
			log.Printf("DB error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
			return
		}

		respondBaitos(w, baitos)
	}
}

func GetBaitoApplicants(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		baitoID := utils.GetParam(r, "baitoid")

		results, err := findBaitoApplicantsFromDB(ctx, app, baitoID)
		if err != nil {
			log.Printf("DB error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, results)
	}
}

func GetMyApplications(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := utils.GetUserIDFromRequest(r)

		results, err := findMyApplicationsFromDB(ctx, app, userID)
		if err != nil {
			logger.Printf("Aggregate error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch applications")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, results)
	}
}
