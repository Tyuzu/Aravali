// File: internal/comments/commentsSQLDB.go

package comments

import (
	"context"
	"strconv"
	"strings"

	"scav/config"
	"scav/infra"
)

var commentsTable = config.Tables.CommentsTable

func insertComment(ctx context.Context, app *infra.Deps, comment Comment) error {
	return nil
}

func findCommentByID(ctx context.Context, app *infra.Deps, commentID string, comment *Comment) error {
	return nil
}

func updateCommentContent(ctx context.Context, app *infra.Deps, commentID string, update map[string]any) (int64, error) {
	return 0, nil
}

func deleteComment(ctx context.Context, app *infra.Deps, commentID, userID string) (int64, error) {
	return 0, nil
}

func findCommentsByEntity(
	ctx context.Context,
	app *infra.Deps,
	entityType string,
	entityID string,
	opts map[string]any,
	comments *[]Comment,
) error {
	return nil
}

func buildCommentListOptions(page, limit int, sortBy string) map[string]any {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	if sortBy == "" {
		sortBy = "created_at DESC"
	}
	return map[string]any{
		"page":   page,
		"limit":  limit,
		"sortBy": sortBy,
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
