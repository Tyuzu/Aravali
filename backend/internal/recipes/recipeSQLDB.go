// File: internal/recipes/recipeSQLDB.go

package recipes

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

// central table name
var recipeTable = config.Tables.RecipeTable

// DB wrappers for recipe operations
func GetRecipeByID(ctx context.Context, app *infra.Deps, id string, out *Recipe) error {
	query := "recipeid = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, recipeTable, query, args, out)
}

func InsertRecipe(ctx context.Context, app *infra.Deps, rec Recipe) error {
	return app.SQLDB.InsertOne(ctx, recipeTable, rec)
}

func UpdateRecipeByID(ctx context.Context, app *infra.Deps, id string, update map[string]any) (int64, error) {
	query := "recipeid = $1"
	args := []any{id}

	return app.SQLDB.Update(ctx, recipeTable, query, args, update)
}

func FindRecipesWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions, out *[]Recipe) error {
	return app.SQLDB.FindManyWithOptions(ctx, recipeTable, query, args, opts, out)
}

func CountRecipes(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.CountDocuments(ctx, recipeTable, query, args)
}

func AggregateRecipes(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.Aggregate(ctx, recipeTable, query, args, out)
}
