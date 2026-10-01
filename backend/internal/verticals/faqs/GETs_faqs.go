// File: internal/verticals/faqs/GETs_faqs.go

package faqs

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"scav/infra"
	"scav/infra/sqldb"
	"scav/utils"
)

/* =========================
   GET SINGLE FAQ
========================= */

func GetFAQ(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		faqID := utils.GetParam(r, "faqid")
		if faqID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid faq ID")
			return
		}

		var faq FAQ
		err := findFAQByID(ctx, app, faqID, &faq)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "FAQ not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, faq)
	}
}

/*
	=========================
	  GET FAQS (PAGINATED)

=========================
*/
func GetFAQs(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		entityType := utils.GetParam(r, "entitytype")
		entityID := utils.GetParam(r, "entityid")

		if entityID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Entity ID is required")
			return
		}

		if !isValidEntityType(entityType) {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid entity type")
			return
		}

		/* ---------- Pagination ---------- */
		page := 1
		limit := 10

		if v := r.URL.Query().Get("page"); v != "" {
			if p, err := strconv.Atoi(v); err == nil && p > 0 {
				page = p
			}
		}

		if v := r.URL.Query().Get("limit"); v != "" {
			if l, err := strconv.Atoi(v); err == nil && l > 0 && l <= 50 {
				limit = l
			}
		}

		skip := (page - 1) * limit
		sortBy := r.URL.Query().Get("sort") // new | old | likes

		/* ---------- Sorting (ORDERED) ---------- */
		orderBy := "created_at DESC, faqid DESC"
		switch sortBy {
		case "old":
			orderBy = "created_at ASC, faqid ASC"
		case "likes":
			orderBy = "likes DESC, created_at DESC"
		}

		opts := sqldb.FindManyOptions{
			Limit:   int64(limit),
			Offset:  int64(skip),
			OrderBy: orderBy,
		}

		var faqs []FAQ
		if err := findFAQsByEntity(ctx, app, entityType, entityID, opts, &faqs); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to fetch faqs"})
			return
		}

		// Always return array (never null)
		if faqs == nil {
			faqs = []FAQ{}
		}

		utils.RespondWithJSON(w, http.StatusOK, faqs)
	}
}
