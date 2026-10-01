// File: internal/products/GETs_products.go

package products

import (
	"context"
	"encoding/json"
	"net/http"
	"scav/infra"
	"scav/infra/sqldb"
	"scav/internal/farms"
	"scav/utils"
	"strings"
	"time"
)

// --------------------------------------------------
// Items
// --------------------------------------------------

func GetItems(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		query := "1 = 1"
		args := []any{}
		if t := r.URL.Query().Get("type"); t != "" {
			query += " AND type = $1"
			args = append(args, t)
		}
		if c := r.URL.Query().Get("category"); c != "" {
			query += " AND category = $2"
			args = append(args, c)
		}
		if s := r.URL.Query().Get("search"); s != "" {
			query += " AND LOWER(name) LIKE $3"
			args = append(args, "%"+strings.ToLower(s)+"%")
		}

		skip, limit := utils.ParsePagination(r, 10, 100)
		orderBy := "name ASC"
		switch r.URL.Query().Get("sort") {
		case "price_asc":
			orderBy = "price ASC"
		case "price_desc":
			orderBy = "price DESC"
		case "name_desc":
			orderBy = "name DESC"
		}

		opts := sqldb.FindManyOptions{
			Offset:  int64(skip),
			Limit:   int64(limit),
			OrderBy: orderBy,
		}

		var items []farms.Product
		if err := FindProductsWithOptions(ctx, app, query, args, opts, &items); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch items")
			return
		}

		total, err := CountProducts(ctx, app, query, args)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to count items")
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
