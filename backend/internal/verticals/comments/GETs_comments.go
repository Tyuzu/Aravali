// File: internal/verticals/comments/GETs_comments.go

package comments

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
   GET SINGLE COMMENT
========================= */

func GetComment(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		commentID := utils.GetParam(r, "commentid")
		if commentID == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid comment ID")
			return
		}

		var comment Comment
		err := findCommentByID(ctx, app, commentID, &comment)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Comment not found")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, comment)
	}
}

/*
	=========================
	  GET COMMENTS (PAGINATED)

=========================
*/
func GetComments(app *infra.Deps) http.HandlerFunc {
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
		orderBy := "created_at DESC, commentid DESC"
		switch sortBy {
		case "old":
			orderBy = "created_at ASC, commentid ASC"
		case "likes":
			orderBy = "likes DESC, created_at DESC"
		}

		opts := sqldb.FindManyOptions{
			Limit:   int64(limit),
			Offset:  int64(skip),
			OrderBy: orderBy,
		}

		var comments []Comment
		if err := findCommentsByEntity(ctx, app, entityType, entityID, opts, &comments); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to fetch comments"})
			return
		}

		// Always return array (never null)
		if comments == nil {
			comments = []Comment{}
		}

		utils.RespondWithJSON(w, http.StatusOK, comments)
	}
}
