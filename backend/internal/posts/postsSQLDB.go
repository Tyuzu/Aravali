package posts

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/db"
)

var blogPostsTable = config.Tables.BlogPostsTable
var usersTable = config.Tables.UserTable

// Wrappers
func SQLGetPostByID(ctx context.Context, app *infra.Deps, id string, out *BlogPost) error {
	return app.DB.FindOne(ctx, blogPostsTable, map[string]any{"postid": id}, out)
}

func SQLFindPostsWithOptions(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions, out *[]BlogPost) error {
	return app.DB.FindManyWithOptions(ctx, blogPostsTable, filter, opts, out)
}

func SQLFindUsersByFilter(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindMany(ctx, usersTable, filter, out)
}

func SQLUpdatePostByFilter(ctx context.Context, app *infra.Deps, filter map[string]any, update any) (any, error) {
	return app.DB.UpdateOne(ctx, blogPostsTable, filter, update)
}

func SQLInsertPost(ctx context.Context, app *infra.Deps, post BlogPost) error {
	return app.DB.InsertOne(ctx, blogPostsTable, post)
}

func SQLDeletePostByFilter(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	return app.DB.DeleteOne(ctx, blogPostsTable, filter)
}

func SQLFindRelatedPostsWithOptions(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions, out any) error {
	return app.DB.FindManyWithOptions(ctx, blogPostsTable, filter, opts, out)
}
