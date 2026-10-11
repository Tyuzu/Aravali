// File: internal/comments/GETs_comments.go

package comments

import (
	"context"
	"net/http"
	"time"

	"scav/infra"
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

		page, limit := parseCommentPageLimit(r.URL.Query().Get("page"), r.URL.Query().Get("limit"))
		sortBy := r.URL.Query().Get("sort")

		comments, err := getCommentPage(ctx, app, entityType, entityID, page, limit, sortBy)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to fetch comments"})
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, comments)
	}
}
