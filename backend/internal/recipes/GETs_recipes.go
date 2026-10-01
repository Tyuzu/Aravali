// File: internal/recipes/GETs_recipes.go

package recipes

import (
	"context"
	"net/http"
	"strings"
	"time"

	"scav/infra"
	"scav/infra/sqldb"
	"scav/utils"
)

// --- Get single recipe ---

func GetRecipe(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := utils.GetParam(r, "id")

		var recipe Recipe
		if err := GetRecipeByID(ctx, app, id, &recipe); err != nil {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}

		normalizeRecipeSlices(&recipe)

		utils.RespondWithJSON(w, http.StatusOK, recipe)
	}
}

// --- List Recipes ---

func GetRecipes(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		query := "1 = 1"
		args := []any{}
		if search := r.URL.Query().Get("search"); search != "" {
			query += " AND (LOWER(title) LIKE $1 OR LOWER(description) LIKE $1)"
			args = append(args, "%"+strings.ToLower(search)+"%")
		}
		if ing := r.URL.Query().Get("ingredient"); ing != "" {
			query += " AND ingredients::text ILIKE $2"
			args = append(args, "%"+strings.ToLower(ing)+"%")
		}
		if tags := r.URL.Query().Get("tags"); tags != "" {
			query += " AND tags::text ILIKE $3"
			args = append(args, "%"+strings.ToLower(tags)+"%")
		}

		skip, limit := utils.ParsePagination(r, 10, 100)
		orderBy := "created_at DESC"
		switch r.URL.Query().Get("sort") {
		case "oldest":
			orderBy = "created_at ASC"
		case "views":
			orderBy = "views DESC"
		case "prepTime":
			orderBy = "preptime ASC"
		}

		opts := sqldb.FindManyOptions{
			Offset:  int64(skip),
			Limit:   int64(limit),
			OrderBy: orderBy,
		}

		var recipes []Recipe
		if err := FindRecipesWithOptions(ctx, app, query, args, opts, &recipes); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch recipes")
			return
		}

		for i := range recipes {
			normalizeRecipeSlices(&recipes[i])
		}

		totalCount, err := CountRecipes(ctx, app, query, args)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to count recipes")
			return
		}

		hasMore := (skip + limit) < int(totalCount)

		utils.RespondWithJSON(w, http.StatusOK, map[string]any{
			"recipes": recipes,
			"hasMore": hasMore,
		})
	}
}

// --- Tags ---

type recipeTagAgg struct {
	Tags []string `db:"tags"`
}

func GetRecipeTags(app *infra.Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		query := "SELECT DISTINCT unnest(tags) AS tags FROM " + recipeTable + " WHERE tags IS NOT NULL"
		var result []recipeTagAgg
		if err := AggregateRecipes(ctx, app, query, nil, &result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tags := []string{}
		for _, item := range result {
			tags = append(tags, item.Tags...)
		}

		utils.RespondWithJSON(w, http.StatusOK, tags)
	}
}
