package reviews

import (
	"context"
	"net/http"
	"scav/infra"
	"scav/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

/* -------------------------
   Get Reviews (list)
------------------------- */

func GetReviews(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		entityType := utils.GetParam(r, "entityType")
		entityId := utils.GetParam(r, "entityId")

		skip, limit := utils.ParsePagination(r, 10, 100)

		var reviews []Review
		if err := GetReviewsForEntity(ctx, app, entityType, entityId, &reviews); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch reviews"})
			return
		}

		utils.SortAndSlice(
			&reviews,
			[]bson.E{{Key: "createdAt", Value: -1}},
			int64(skip),
			int64(limit),
		)

		if reviews == nil {
			reviews = []Review{}
		}

		utils.RespondWithJSON(w, http.StatusOK, reviews)
	}
}

/* -------------------------
   Get Review (single)
------------------------- */

func GetReview(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reviewId := utils.GetParam(r, "reviewId")

		var review Review
		if err := GetReviewByID(r.Context(), app, reviewId, &review); err != nil {
			utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Review not found"})
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, review)
	}
}
