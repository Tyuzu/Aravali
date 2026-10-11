// File: internal/verticals/comments/commentsSQLDB.go

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
}

func findCommentByID(ctx context.Context, app *infra.Deps, commentID string, comment *Comment) error {
}

func updateCommentContent(ctx context.Context, app *infra.Deps, commentID string, update map[string]any) (int64, error) {
}

func deleteComment(ctx context.Context, app *infra.Deps, commentID, userID string) (int64, error) {
}

func findCommentsByEntity(
	ctx context.Context,
	app *infra.Deps,
	entityType string,
	entityID string,
	opts map[string]any,
	comments *[]Comment,
) error {
}

func buildCommentListOptions(page, limit int, sortBy string) map[string]any {

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
