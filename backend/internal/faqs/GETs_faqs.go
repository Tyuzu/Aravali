// File: internal/faqs/GETs_faqs.go

package faqs

import (
	"context"
	"net/http"
	"time"

	"scav/infra"
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

		page, limit := parseFAQPageLimit(r.URL.Query().Get("page"), r.URL.Query().Get("limit"))
		sortBy := r.URL.Query().Get("sort")

		faqs, err := getFAQPage(ctx, app, entityType, entityID, page, limit, sortBy)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to fetch faqs"})
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, faqs)
	}
}
