// File: internal/posts/postsSQLDB.go

package posts

import (
	"context"
	"strconv"

	"scav/config"
	"scav/infra"
)

var blogPostsTable = config.Tables.BlogPostsTable
var usersTable = config.Tables.UserTable

// Wrappers
func GetPostByID(ctx context.Context, app *infra.Deps, id string, out *BlogPost) error {
}

func FindPostsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any, out *[]BlogPost) error {
}

func FindUsersByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
}

func UpdatePostByFilter(ctx context.Context, app *infra.Deps, query string, args []any, update map[string]any) (int64, error) {
}

func InsertPost(ctx context.Context, app *infra.Deps, post BlogPost) error {
}

func DeletePostByFilter(ctx context.Context, app *infra.Deps, query string, args []any) (int64, error) {
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

func buildPostListOptions(limit, page int) map[string]any {
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

func FindRelatedPostsWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any, out any) error {
}
