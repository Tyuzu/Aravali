// File: internal/recipes/recipeSQLDB.go

package recipes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"scav/config"
	"scav/infra"
)

// central table name
var recipeTable = config.Tables.RecipeTable

// DB wrappers for recipe operations
func GetRecipeByID(ctx context.Context, app *infra.Deps, id string, out *Recipe) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if out == nil {
		return errors.New("recipe output is nil")
	}
	row := app.SQLDB.QueryRow(ctx, `SELECT recipeid, userid, title, description, cookTime, cuisine, dietary, portionSize, season, tags, images, ingredients, steps, difficulty, banner, servings, videoUrl, notes, createdAt, views FROM recipes WHERE recipeid = $1 LIMIT 1`, id)
	var recipeID, userID, title, description, cookTime, cuisine, portionSize, season, difficulty, banner, videoURL, notes string
	var dietary, tags, images, steps []string
	var ingredients []byte
	var servings int
	var createdAt int64
	var views int
	if err := row.Scan(&recipeID, &userID, &title, &description, &cookTime, &cuisine, &dietary, &portionSize, &season, &tags, &images, &ingredients, &steps, &difficulty, &banner, &servings, &videoURL, &notes, &createdAt, &views); err != nil {
		return err
	}
	recipe := &Recipe{RecipeId: recipeID, UserID: userID, Title: title, Description: description, CookTime: cookTime, Cuisine: cuisine, Dietary: dietary, PortionSize: portionSize, Season: season, Tags: tags, Images: images, Difficulty: difficulty, Banner: banner, Servings: servings, VideoURL: videoURL, Notes: notes, CreatedAt: createdAt, Views: views}
	if len(ingredients) > 0 && string(ingredients) != "null" {
		if err := json.Unmarshal(ingredients, &recipe.Ingredients); err != nil {
			return err
		}
	}
	normalizeRecipeSlices(recipe)
	*out = *recipe
	return nil
}

func InsertRecipe(ctx context.Context, app *infra.Deps, rec Recipe) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if strings.TrimSpace(rec.RecipeId) == "" {
		return errors.New("recipe ID required")
	}
	normalizeRecipeSlices(&rec)
	ingredients, err := json.Marshal(rec.Ingredients)
	if err != nil {
		return err
	}
	_, err = app.SQLDB.Exec(ctx,
		`INSERT INTO recipes (recipeid, userid, title, description, cookTime, cuisine, dietary, portionSize, season, tags, images, ingredients, steps, difficulty, banner, servings, videoUrl, notes, createdAt, views, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, NOW(), NOW(), '{}'::jsonb)
		ON CONFLICT (recipeid) DO UPDATE SET
			userid = EXCLUDED.userid,
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			cookTime = EXCLUDED.cookTime,
			cuisine = EXCLUDED.cuisine,
			dietary = EXCLUDED.dietary,
			portionSize = EXCLUDED.portionSize,
			season = EXCLUDED.season,
			tags = EXCLUDED.tags,
			images = EXCLUDED.images,
			ingredients = EXCLUDED.ingredients,
			steps = EXCLUDED.steps,
			difficulty = EXCLUDED.difficulty,
			banner = EXCLUDED.banner,
			servings = EXCLUDED.servings,
			videoUrl = EXCLUDED.videoUrl,
			notes = EXCLUDED.notes,
			createdAt = EXCLUDED.createdAt,
			views = EXCLUDED.views,
			updated_at = NOW()`,
		rec.RecipeId,
		rec.UserID,
		rec.Title,
		rec.Description,
		rec.CookTime,
		rec.Cuisine,
		rec.Dietary,
		rec.PortionSize,
		rec.Season,
		rec.Tags,
		rec.Images,
		ingredients,
		rec.Steps,
		rec.Difficulty,
		rec.Banner,
		rec.Servings,
		rec.VideoURL,
		rec.Notes,
		rec.CreatedAt,
		rec.Views,
	)
	return err
}

func UpdateRecipeByID(ctx context.Context, app *infra.Deps, id string, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if len(update) == 0 {
		return 0, nil
	}
	parts := make([]string, 0, len(update))
	args := make([]any, 0, len(update)+1)
	for key, value := range update {
		parts = append(parts, key+" = $"+fmt.Sprintf("%d", len(args)+1))
		args = append(args, value)
	}
	parts = append(parts, "updated_at = NOW()")
	query := `UPDATE recipes SET ` + strings.Join(parts, ", ") + ` WHERE recipeid = $` + fmt.Sprintf("%d", len(args)+1)
	args = append(args, id)
	result, err := app.SQLDB.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func FindRecipesWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any, out *[]Recipe) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if out == nil {
		return errors.New("recipe output is nil")
	}
	limit := 25
	offset := 0
	if opts != nil {
		if v, ok := opts["limit"].(int); ok {
			limit = v
		}
		if v, ok := opts["offset"].(int); ok {
			offset = v
		}
		if v, ok := opts["sort"].(string); ok && v != "" {
			_ = v
		}
	}
	if limit <= 0 {
		limit = 25
	}
	// Use a simple, safe ordering; query is already filtered before this call.
	q := fmt.Sprintf(`SELECT recipeid, userid, title, description, cookTime, cuisine, dietary, portionSize, season, tags, images, ingredients, steps, difficulty, banner, servings, videoUrl, notes, createdAt, views FROM recipes WHERE %s ORDER BY createdAt DESC LIMIT %d OFFSET %d`, query, limit, offset)
	rows, err := app.SQLDB.Query(ctx, q, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	items := make([]Recipe, 0)
	for rows.Next() {
		var recipeID, userID, title, description, cookTime, cuisine, portionSize, season, difficulty, banner, videoURL, notes string
		var dietary, tags, images, steps []string
		var ingredients []byte
		var servings int
		var createdAt int64
		var views int
		if err := rows.Scan(&recipeID, &userID, &title, &description, &cookTime, &cuisine, &dietary, &portionSize, &season, &tags, &images, &ingredients, &steps, &difficulty, &banner, &servings, &videoURL, &notes, &createdAt, &views); err != nil {
			return err
		}
		recipe := Recipe{RecipeId: recipeID, UserID: userID, Title: title, Description: description, CookTime: cookTime, Cuisine: cuisine, Dietary: dietary, PortionSize: portionSize, Season: season, Tags: tags, Images: images, Difficulty: difficulty, Banner: banner, Servings: servings, VideoURL: videoURL, Notes: notes, CreatedAt: createdAt, Views: views}
		if len(ingredients) > 0 && string(ingredients) != "null" {
			if err := json.Unmarshal(ingredients, &recipe.Ingredients); err != nil {
				return err
			}
		}
		normalizeRecipeSlices(&recipe)
		items = append(items, recipe)
	}
	*out = items
	return rows.Err()
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

func buildRecipeListOptions(skip, limit int, sort string) map[string]any {
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
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	var total int64
	if err := app.SQLDB.QueryRow(ctx, `SELECT COUNT(*) FROM recipes WHERE `+query, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func AggregateRecipes(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if out == nil {
		return errors.New("aggregate output is nil")
	}
	rows, err := app.SQLDB.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return errors.New("aggregate output must be pointer to slice")
	}
	itemType := v.Elem().Type().Elem()
	for rows.Next() {
		if itemType.Kind() == reflect.Struct {
			item := reflect.New(itemType).Elem()
			if item.NumField() == 1 {
				var tag string
				if err := rows.Scan(&tag); err != nil {
					return err
				}
				field := item.Field(0)
				if field.CanSet() && field.Kind() == reflect.Slice {
					field.Set(reflect.ValueOf([]string{tag}))
				}
				v.Elem().Set(reflect.Append(v.Elem(), item))
				continue
			}
		}
		values, err := rows.Values()
		if err != nil {
			return err
		}
		item := reflect.New(itemType).Elem()
		for i := 0; i < item.NumField() && i < len(values); i++ {
			field := item.Field(i)
			if !field.CanSet() {
				continue
			}
			if values[i] == nil {
				continue
			}
			if field.Type() == reflect.TypeOf(values[i]) {
				field.Set(reflect.ValueOf(values[i]))
			}
		}
		v.Elem().Set(reflect.Append(v.Elem(), item))
	}
	return rows.Err()
}
