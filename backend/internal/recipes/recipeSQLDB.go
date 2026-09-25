package recipes

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/db"
)

// central table name
var recipeTable = config.Tables.RecipeTable

// DB wrappers for recipe operations
func SQLGetRecipeByID(ctx context.Context, app *infra.Deps, id string, out *Recipe) error {
	return app.DB.FindOne(ctx, recipeTable, map[string]any{"recipeid": id}, out)
}

func SQLInsertRecipe(ctx context.Context, app *infra.Deps, rec Recipe) error {
	return app.DB.InsertOne(ctx, recipeTable, rec)
}

func SQLUpdateRecipeByID(ctx context.Context, app *infra.Deps, id string, update any) (any, error) {
	return app.DB.Update(ctx, recipeTable, map[string]any{"recipeid": id}, update)
}

func SQLFindRecipesWithOptions(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions, out *[]Recipe) error {
	return app.DB.FindManyWithOptions(ctx, recipeTable, filter, opts, out)
}

func SQLCountRecipes(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	return app.DB.CountDocuments(ctx, recipeTable, filter)
}

func SQLAggregateRecipes(ctx context.Context, app *infra.Deps, pipeline any, out any) error {
	return app.DB.Aggregate(ctx, recipeTable, pipeline, out)
}
