package posts

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var blogPostsTable = config.Tables.BlogPostsTable
var usersTable = config.Tables.UserTable

// Wrappers
func SQLGetPostByID(ctx context.Context, app *infra.Deps, id string, out *BlogPost) error {
	query := "postid = $1"
	args := []any{id}

	return app.SQLDB.FindOne(ctx, blogPostsTable, query, args, out)
}

func SQLFindPostsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions, out *[]BlogPost) error {
	return app.SQLDB.FindManyWithOptions(ctx, blogPostsTable, query, args, opts, out)
}

func SQLFindUsersByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindMany(ctx, usersTable, query, args, out)
}

func SQLUpdatePostByFilter(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
	return app.SQLDB.UpdateOne(ctx, blogPostsTable, query, args, update)
}

func SQLInsertPost(ctx context.Context, app *infra.Deps, post BlogPost) error {
	return app.SQLDB.InsertOne(ctx, blogPostsTable, post)
}

func SQLDeletePostByFilter(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
	return app.SQLDB.DeleteOne(ctx, blogPostsTable, query, args)
}

func SQLFindRelatedPostsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions, out any) error {
	return app.SQLDB.FindManyWithOptions(ctx, blogPostsTable, query, args, opts, out)
}
