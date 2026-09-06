package posts

import (
	"context"

	"scav/config"
	"scav/infra"
	"scav/infra/db"
)

var blogPostsCollection = config.Collections.BlogPostsCollection
var usersCollection = config.Collections.UserCollection

// Wrappers
func GetPostByID(ctx context.Context, app *infra.Deps, id string, out *BlogPost) error {
	return app.DB.FindOne(ctx, blogPostsCollection, map[string]any{"postid": id}, out)
}

func FindPostsWithOptions(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions, out *[]BlogPost) error {
	return app.DB.FindManyWithOptions(ctx, blogPostsCollection, filter, opts, out)
}

func FindUsersByFilter(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindMany(ctx, usersCollection, filter, out)
}

func UpdatePostByFilter(ctx context.Context, app *infra.Deps, filter map[string]any, update any) (any, error) {
	return app.DB.UpdateOne(ctx, blogPostsCollection, filter, update)
}

func InsertPost(ctx context.Context, app *infra.Deps, post BlogPost) error {
	return app.DB.InsertOne(ctx, blogPostsCollection, post)
}

func DeletePostByFilter(ctx context.Context, app *infra.Deps, filter map[string]any) (int64, error) {
	return app.DB.DeleteOne(ctx, blogPostsCollection, filter)
}

func FindRelatedPostsWithOptions(ctx context.Context, app *infra.Deps, filter map[string]any, opts db.FindManyOptions, out any) error {
	return app.DB.FindManyWithOptions(ctx, blogPostsCollection, filter, opts, out)
}
