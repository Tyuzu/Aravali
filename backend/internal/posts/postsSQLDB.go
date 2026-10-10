// File: internal/posts/postsSQLDB.go

package posts

import (
	"context"
	"strconv"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var blogPostsTable = config.Tables.BlogPostsTable
var usersTable = config.Tables.UserTable

// Wrappers
func GetPostByID(ctx context.Context, app *infra.Deps, id string, out *BlogPost) error {
	query := "postid = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, blogPostsTable, query, args, out)
}

func FindPostsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions, out *[]BlogPost) error {
	return app.SQLDB.FindManyWithOptions(ctx, blogPostsTable, query, args, opts, out)
}

func FindUsersByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindMany(ctx, usersTable, query, args, out)
}

func UpdatePostByFilter(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	return app.SQLDB.UpdateOne(ctx, blogPostsTable, query, args, update)
}

func InsertPost(ctx context.Context, app *infra.Deps, post BlogPost) error {
	return app.SQLDB.InsertOne(ctx, blogPostsTable, post)
}

func DeletePostByFilter(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.DeleteOne(ctx, blogPostsTable, query, args)
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

func buildRelatedPostsOptions() sqldb.FindManyOptions {
	return sqldb.FindManyOptions{
		Limit:   10,
		OrderBy: "created_at DESC",
	}
}

func GetRelatedPostsPage(ctx context.Context, app *infra.Deps, postID, category, subcategory string, tags []string) ([]Post, error) {
	query, args := buildRelatedPostsQuery(postID, category, subcategory, tags)
	opts := buildRelatedPostsOptions()

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

func buildPostListOptions(limit, page int) sqldb.FindManyOptions {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if page < 1 {
		page = 1
	}

	skipp := (page - 1) * limit
	return sqldb.FindManyOptions{
		Limit:   int64(limit),
		Offset:  int64(skipp),
		OrderBy: "created_at DESC",
	}
}

func GetPostsPage(ctx context.Context, app *infra.Deps, limit, page int) ([]BlogPost, error) {
	query, args := buildPostListQuery()
	opts := buildPostListOptions(limit, page)

	var posts []BlogPost
	if err := FindPostsWithOptions(ctx, app, query, args, opts, &posts); err != nil {
		return nil, err
	}
	if posts == nil {
		return []BlogPost{}, nil
	}
	return posts, nil
}

func FindRelatedPostsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions, out any) error {
	return app.SQLDB.FindManyWithOptions(ctx, blogPostsTable, query, args, opts, out)
}
