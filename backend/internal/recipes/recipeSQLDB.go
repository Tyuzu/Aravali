// File: internal/recipes/recipeSQLDB.go

package recipes

import (
	"context"
	"strings"

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

func buildRecipeSearchQuery(search, ingredient, tags string) (string, []any) {
	query := "1 = 1"
	args := []any{}

	if search != "" {
		query += " AND (LOWER(title) LIKE $1 OR LOWER(description) LIKE $1)"
		args = append(args, "%"+strings.ToLower(search)+"%")
	}
	if ingredient != "" {
		query += " AND ingredients::text ILIKE $2"
		args = append(args, "%"+strings.ToLower(ingredient)+"%")
	}
	if tags != "" {
		query += " AND tags::text ILIKE $3"
		args = append(args, "%"+strings.ToLower(tags)+"%")
	}

	return query, args
}

func buildRecipeListOptions(skip, limit int, sort string) sqldb.FindManyOptions {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if skip < 0 {
		skip = 0
	}

	orderBy := "created_at DESC"
	switch sort {
	case "oldest":
		orderBy = "created_at ASC"
	case "views":
		orderBy = "views DESC"
	case "prepTime":
		orderBy = "preptime ASC"
	}

	return sqldb.FindManyOptions{
		Offset:  int64(skip),
		Limit:   int64(limit),
		OrderBy: orderBy,
	}
}

func getRecipesPage(ctx context.Context, app *infra.Deps, search, ingredient, tags string, skip, limit int, sort string) ([]Recipe, int64, error) {
	query, args := buildRecipeSearchQuery(search, ingredient, tags)
	opts := buildRecipeListOptions(skip, limit, sort)

	var recipes []Recipe
	if err := FindRecipesWithOptions(ctx, app, query, args, opts, &recipes); err != nil {
		return nil, 0, err
	}

	total, err := CountRecipes(ctx, app, query, args)
	if err != nil {
		return nil, 0, err
	}

	return recipes, total, nil
}

func CountRecipes(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.CountDocuments(ctx, recipeTable, query, args)
}

func AggregateRecipes(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.Aggregate(ctx, recipeTable, query, args, out)
}
