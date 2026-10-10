// File: internal/products/GETs_products.go

package products

import (
	"context"
	"encoding/json"
	"net/http"
	"scav/infra"
	"scav/utils"
	"time"
)

// --------------------------------------------------
// Items
// --------------------------------------------------

func GetItems(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		itemType := r.URL.Query().Get("type")
		category := r.URL.Query().Get("category")
		search := r.URL.Query().Get("search")
		skip, limit := utils.ParsePagination(r, 10, 100)

		items, total, err := getProductsPage(ctx, app, itemType, category, search, skip, limit, r.URL.Query().Get("sort"))
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch items")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, utils.M{
			"success": true,
			"items":   items,
			"total":   total,
			"page":    skip/limit + 1,
			"limit":   limit,
		})
	}
}

// --------------------------------------------------
// Item Categories
// --------------------------------------------------

func GetItemCategories(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemType := r.URL.Query().Get("type")

		var categories []string

		switch itemType {
		case "tool":
			categories = []string{
				"Cutting Tools",
				"Irrigation Tools",
				"Harvesting Tools",
				"Hand Tools",
				"Protective Gear",
				"Fertilizer Applicators",
			}
		default:
			categories = []string{
				"Spices",
				"Pickles",
				"Flour",
				"Oils",
				"Honey",
				"Tea & Coffee",
				"Dry Fruits",
				"Natural Sweeteners",
			}
		}

		if err := json.NewEncoder(w).Encode(categories); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
		}
	}
}
