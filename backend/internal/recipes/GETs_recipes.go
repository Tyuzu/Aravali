// File: internal/recipes/GETs_recipes.go

package recipes

import (
	"context"
	"net/http"
	"time"

	"scav/infra"
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

		search := r.URL.Query().Get("search")
		ingredient := r.URL.Query().Get("ingredient")
		tags := r.URL.Query().Get("tags")
		skip, limit := utils.ParsePagination(r, 10, 100)

		recipes, totalCount, err := getRecipesPage(ctx, app, search, ingredient, tags, skip, limit, r.URL.Query().Get("sort"))
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch recipes")
			return
		}

		for i := range recipes {
			normalizeRecipeSlices(&recipes[i])
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
