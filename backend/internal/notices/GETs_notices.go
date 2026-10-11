// File: internal/notices/GETs_notices.go

package notices

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"scav/infra"
	"scav/utils"
)

// --- Get notices list (optimized: summary only) ---
func GetNotices(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		entityType := utils.GetParam(r, "entitytype")
		entityID := utils.GetParam(r, "entityid")

		page := 1
		limit := 10
		if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
			page = p
		}
		if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
			limit = l
		}

		sortBy := r.URL.Query().Get("sort")
		notices, err := getNoticesPage(ctx, app, entityType, entityID, page, limit, sortBy)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch notices")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, buildNoticeSummary(notices))
	}
}

// --- Get single notice ---
func GetNotice(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		noticeID := strings.TrimSpace(utils.GetParam(r, "noticeid"))
		if noticeID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		notice, err := findNoticeByID(ctx, app, noticeID)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Notice not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, notice)
	}
}
