package recipes

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/db"
)

// central collection name
var recipeCollection = config.Collections.RecipeCollection

// DB wrappers for recipe operations
func GetRecipeByID(ctx context.Context, app *infra.Deps, id string, out *Recipe) error {
	return app.DB.FindOne(ctx, recipeCollection, map[string]any{"recipeid": id}, out)
}

func InsertRecipe(ctx context.Context, app *infra.Deps, rec Recipe) error {
	return app.DB.InsertOne(ctx, recipeCollection, rec)
}

func UpdateRecipeByID(ctx context.Context, app *infra.Deps, id string, update any) (any, error) {
	return app.DB.Update(ctx, recipeCollection, map[string]any{"recipeid": id}, update)
}

func FindRecipesWithOptions(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions, out *[]Recipe) error {
	return app.DB.FindManyWithOptions(ctx, recipeCollection, filter, opts, out)
}

func CountRecipes(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	return app.DB.CountDocuments(ctx, recipeCollection, filter)
}

func AggregateRecipes(ctx context.Context, app *infra.Deps, pipeline any, out any) error {
	return app.DB.Aggregate(ctx, recipeCollection, pipeline, out)
}
