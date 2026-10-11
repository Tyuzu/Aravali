// File: internal/posts/postsSQLDB.go

package posts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"scav/config"
	"scav/infra"
)

var blogPostsTable = config.Tables.BlogPostsTable
var usersTable = config.Tables.UserTable

// Wrappers
func GetPostByID(ctx context.Context, app *infra.Deps, id string, out *BlogPost) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if id == "" {
		return nil
	}
	row := app.SQLDB.QueryRow(ctx, `SELECT * FROM `+blogPostsTable+` WHERE postid = $1 LIMIT 1`, id)
	if row == nil {
		return nil
	}
	return nil
}

func FindPostsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any, out *[]BlogPost) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}
	_, err := app.SQLDB.Query(ctx, `SELECT * FROM `+blogPostsTable+` WHERE `+query, args...)
	return err
}

func FindUsersByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	_, err := app.SQLDB.Query(ctx, `SELECT * FROM `+usersTable+` WHERE `+query, args...)
	return err
}

func UpdatePostByFilter(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if len(update) == 0 {
		return 0, nil
	}
	setParts := make([]string, 0, len(update))
	setArgs := make([]any, 0, len(update))
	for key, value := range update {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", key, len(setArgs)+1))
		setArgs = append(setArgs, value)
	}
	vals := append(setArgs, args...)
	result, err := app.SQLDB.Exec(ctx, `UPDATE `+blogPostsTable+` SET `+strings.Join(setParts, ", ")+` WHERE `+query, vals...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func InsertPost(ctx context.Context, app *infra.Deps, post BlogPost) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(ctx, `INSERT INTO `+blogPostsTable+` DEFAULT VALUES`)
	return err
}

func DeletePostByFilter(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if query == "" {
		query = "1 = 1"
	}
	result, err := app.SQLDB.Exec(ctx, `DELETE FROM `+blogPostsTable+` WHERE `+query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func buildRelatedPostsQuery(postID, category, subcategory string, tags []string) (string, []any) {
	args := []any{postID}
	query := "postid <> $1"

	if category != "" {
		query += " AND category = $" + strconv.Itoa(len(args)+1)
		args = append(args, category)
	}
	if subcategory != "" {
		query += " AND subcategory = $" + strconv.Itoa(len(args)+1)
		args = append(args, subcategory)
	}
	if len(tags) > 0 {
		query += " AND tags && $" + strconv.Itoa(len(args)+1)
		args = append(args, tags)
	}

	return query, args
}

func buildRelatedPostsOptions() map[string]any {
	return map[string]any{"limit": 10, "offset": 0}
}

func GetRelatedPostsPage(ctx context.Context, app *infra.Deps, postID, category, subcategory string, tags []string) ([]Post, error) {
	query, args := buildRelatedPostsQuery(postID, category, subcategory, tags)
	opts := buildRelatedPostsOptions()
	_ = opts
	var related []Post
	if err := FindRelatedPostsWithOptions(ctx, app, query, args, opts, &related); err != nil {
		return nil, err
	}
	if related == nil {
		return []Post{}, nil
	}
	return related, nil
}

func buildPostListQuery() (string, []any) {
	return "1 = 1", nil
}

func buildPostListOptions(limit, page int) map[string]any {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	return map[string]any{"limit": limit, "offset": (page - 1) * limit}
}

func GetPostsPage(ctx context.Context, app *infra.Deps, limit, page int) ([]BlogPost, error) {
	query, args := buildPostListQuery()
	opts := buildPostListOptions(limit, page)
	_ = opts
	var posts []BlogPost
	if err := FindPostsWithOptions(ctx, app, query, args, opts, &posts); err != nil {
		return nil, err
	}
	if posts == nil {
		return []BlogPost{}, nil
	}
	return posts, nil
}

func FindRelatedPostsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any, out any) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}
	_, err := app.SQLDB.Query(ctx, `SELECT * FROM `+blogPostsTable+` WHERE `+query, args...)
	return err
}
