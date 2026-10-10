// File: internal/verticals/comments/commentsSQLDB.go

package comments

import (
	"context"
	"strconv"
	"strings"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var commentsTable = config.Tables.CommentsTable

func insertComment(ctx context.Context, app *infra.Deps, comment Comment) error {
	return app.SQLDB.Insert(ctx, commentsTable, comment)
}

func findCommentByID(ctx context.Context, app *infra.Deps, commentID string, comment *Comment) error {
	query := "commentid = $1"
	args := []any{commentID}

	return app.SQLDB.FindOne(ctx, commentsTable, query, args, comment)
}

func updateCommentContent(ctx context.Context, app *infra.Deps, commentID string, update map[string]any) (int64, error) {
	query := "commentid = $1"
	args := []any{commentID}

	return app.SQLDB.UpdateOne(ctx, commentsTable, query, args, update)
}

func deleteComment(ctx context.Context, app *infra.Deps, commentID, userID string) (int64, error) {
	query := "commentid = $1 AND createdby = $2"
	args := []any{commentID, userID}

	return app.SQLDB.Delete(ctx, commentsTable, query, args)
}

func findCommentsByEntity(
	ctx context.Context,
	app *infra.Deps,
	entityType string,
	entityID string,
	opts sqldb.FindManyOptions,
	comments *[]Comment,
) error {
	query := "entity_type = $1 AND entity_id = $2"
	args := []any{entityType, entityID}

	return app.SQLDB.FindManyWithOptions(ctx, commentsTable, query, args, opts, comments)
}

func buildCommentListOptions(page, limit int, sortBy string) sqldb.FindManyOptions {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	skip := (page - 1) * limit
	orderBy := "created_at DESC, commentid DESC"
	switch sortBy {
	case "old":
		orderBy = "created_at ASC, commentid ASC"
	case "likes":
		orderBy = "likes DESC, created_at DESC"
	}

	return sqldb.FindManyOptions{
		Limit:   int64(limit),
		Offset:  int64(skip),
		OrderBy: orderBy,
	}
}

func getCommentPage(ctx context.Context, app *infra.Deps, entityType, entityID string, page, limit int, sortBy string) ([]Comment, error) {
	var comments []Comment
	opts := buildCommentListOptions(page, limit, sortBy)
	if err := findCommentsByEntity(ctx, app, entityType, entityID, opts, &comments); err != nil {
		return nil, err
	}
	if comments == nil {
		return []Comment{}, nil
	}
	return comments, nil
}

func parseCommentPageLimit(rawPage, rawLimit string) (int, int) {
	page := 1
	if v := strings.TrimSpace(rawPage); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if v := strings.TrimSpace(rawLimit); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	return page, limit
}
