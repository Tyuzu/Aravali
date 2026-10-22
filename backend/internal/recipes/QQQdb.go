// File: internal/recipes/recipeSQLDB.go

package recipes

import (
	"context"

	"scav/config"
	"scav/infra"
)

// central table name
var recipeTable = config.Tables.RecipeTable

// DB wrappers for recipe operations
func QQQGetRecipeByID(ctx context.Context, app *infra.Deps, id string, out *Recipe) error {
	return ctx.Err()
}

func QQQInsertRecipe(ctx context.Context, app *infra.Deps, rec Recipe) error {

	return ctx.Err()
}

func QQQUpdateRecipeByID(ctx context.Context, app *infra.Deps, id string, update map[string]any) (int64, error) {
	return 1, nil
}

func QQQFindRecipesWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any, out *[]Recipe) error {
	return ctx.Err()
}

func QQQbuildRecipeSearchQuery(search, ingredient, tags string) (string, []any) {
	query := "1 = 1"
	args := []any{}

	return query, args
}

func QQQbuildRecipeListOptions(skip, limit int, sort string) map[string]any {
	if limit <= 0 {
		limit = 25
	}
	if skip < 0 {
		skip = 0
	}
	if sort == "" {
		sort = "createdAt DESC"
	}
	return map[string]any{"offset": skip, "limit": limit, "sort": sort}
}

func QQQgetRecipesPage(ctx context.Context, app *infra.Deps, search, ingredient, tags string, skip, limit int, sort string) ([]Recipe, int64, error) {

	var total int64
	var recipes []Recipe
	return recipes, total, nil
}

func QQQCountRecipes(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {

	var total int64
	return total, nil
}

func QQQAggregateRecipes(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return ctx.Err()
}
