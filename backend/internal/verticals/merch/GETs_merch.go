package merch

import (
	"context"
	"net/http"
	"time"

	"scav/infra"
	"scav/utils"
)

func GetMerch(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entityType := utils.GetParam(r, "entityType")
		eventID := utils.GetParam(r, "eventid")
		merchID := utils.GetParam(r, "merchid")

		if !validateEntityType(entityType) {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   "invalid entity type",
			})
			return
		}

		merch, err := findMerchByEntity(r.Context(), app, entityType, eventID, merchID)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusNotFound, map[string]any{
				"success": false,
				"error":   "merch not found",
			})
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    merch,
		})
	}
}

func GetMerchs(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entityType := utils.GetParam(r, "entityType")
		eventID := utils.GetParam(r, "eventid")

		if !validateEntityType(entityType) {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]any{
				"success": false,
				"error":   "invalid entity type",
			})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		list, err := findMerchsByEntity(ctx, app, entityType, eventID)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]any{
				"success": false,
				"error":   "failed to fetch merch",
			})
			return
		}

		if list == nil {
			list = []Merch{}
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    list,
		})
	}
}

func GetMerchPage(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		merchID := utils.GetParam(r, "entityType") // route constraint

		merch, err := findMerchByMerchID(r.Context(), app, merchID)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusNotFound, map[string]any{
				"success": false,
				"error":   "merch not found",
			})
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    merch,
		})
	}
}
