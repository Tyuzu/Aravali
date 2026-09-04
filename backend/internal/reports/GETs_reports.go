package reports

import (
	"net/http"
	"strconv"
	"strings"

	"scav/infra"
	"scav/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAppeals(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.URL.Query()

		status := strings.TrimSpace(q.Get("status"))
		if status == "" {
			status = "pending"
		}

		filter := map[string]any{"status": status}

		limit := int64(20)
		offset := int64(0)

		if l := strings.TrimSpace(q.Get("limit")); l != "" {
			if v, err := strconv.ParseInt(l, 10, 64); err == nil && v > 0 {
				limit = v
			}
		}
		if o := strings.TrimSpace(q.Get("offset")); o != "" {
			if v, err := strconv.ParseInt(o, 10, 64); err == nil && v >= 0 {
				offset = v
			}
		}

		var appeals []map[string]any
		if err := app.DB.FindMany(ctx, appealsCollection, filter, &appeals); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch appeals",
			})
			return
		}

		utils.SortAndSlice(
			&appeals,
			[]bson.E{{Key: "createdAt", Value: -1}},
			offset,
			limit,
		)

		if appeals == nil {
			appeals = []map[string]any{}
		}

		utils.RespondWithJSON(w, http.StatusOK, appeals)
	}
}

func GetMyAppeals(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := utils.GetUserIDFromRequest(r)
		if userID == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		filter := map[string]any{"userid": userID}
		if status := strings.TrimSpace(r.URL.Query().Get("status")); status != "" {
			filter["status"] = status
		}

		var appeals []map[string]any
		if err := app.DB.FindMany(ctx, appealsCollection, filter, &appeals); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch your appeals",
			})
			return
		}
		if appeals == nil {
			appeals = []map[string]any{}
		}

		utils.SortAndSlice(&appeals, []bson.E{{Key: "createdAt", Value: -1}}, 0, int64(len(appeals)))
		utils.RespondWithJSON(w, http.StatusOK, appeals)
	}
}
