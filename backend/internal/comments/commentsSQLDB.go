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
	if app == nil || app.SQLDB == nil {
		return nil
	}
	_, err := app.SQLDB.Exec(
		ctx,
		`INSERT INTO comments (commentid, entity_type, entity_id, userid, content, parent_id, likes, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (commentid) DO UPDATE SET entity_type = EXCLUDED.entity_type, entity_id = EXCLUDED.entity_id, userid = EXCLUDED.userid, content = EXCLUDED.content, parent_id = EXCLUDED.parent_id, likes = EXCLUDED.likes, updated_at = NOW()`,
		comment.CommentID,
		comment.EntityType,
		comment.EntityID,
		comment.CreatedBy,
		comment.Content,
		nil,
		comment.Likes,
		comment.CreatedAt,
		comment.UpdatedAt,
	)
	return err
}

func findCommentByID(ctx context.Context, app *infra.Deps, commentID string, comment *Comment) error {
	if app == nil || app.SQLDB == nil || comment == nil {
		return nil
	}
	var parentID *string
	row := app.SQLDB.QueryRow(ctx, `SELECT commentid, entity_type, entity_id, userid, content, parent_id, likes, created_at, updated_at FROM comments WHERE commentid = $1 LIMIT 1`, commentID)
	if err := row.Scan(&comment.CommentID, &comment.EntityType, &comment.EntityID, &comment.CreatedBy, &comment.Content, &parentID, &comment.Likes, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
		return err
	}
	_ = parentID
	return nil
}

func updateCommentContent(ctx context.Context, app *infra.Deps, commentID string, update map[string]any) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	content, _ := update["content"].(string)
	if strings.TrimSpace(commentID) == "" || content == "" {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(ctx, `UPDATE comments SET content = $1, updated_at = NOW() WHERE commentid = $2`, content, commentID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func deleteComment(ctx context.Context, app *infra.Deps, commentID, userID string) (int64, error) {
	if app == nil || app.SQLDB == nil {
		return 0, nil
	}
	if strings.TrimSpace(commentID) == "" {
		return 0, nil
	}
	result, err := app.SQLDB.Exec(ctx, `DELETE FROM comments WHERE commentid = $1 AND userid = $2`, commentID, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func findCommentsByEntity(
	ctx context.Context,
	app *infra.Deps,
	entityType string,
	entityID string,
	opts map[string]any,
	comments *[]Comment,
) error {
	if app == nil || app.SQLDB == nil || comments == nil {
		return nil
	}
	query := `SELECT commentid, entity_type, entity_id, userid, content, parent_id, likes, created_at, updated_at FROM comments WHERE entity_type = $1 AND entity_id = $2 ORDER BY created_at DESC`
	args := []any{entityType, entityID}
	if limit, ok := opts["limit"].(int); ok && limit > 0 {
		query += ` LIMIT $3`
		args = append(args, limit)
	}
	rows, err := app.SQLDB.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	items := make([]Comment, 0)
	for rows.Next() {
		var item Comment
		var parentID *string
		var likes int64
		if err := rows.Scan(&item.CommentID, &item.EntityType, &item.EntityID, &item.CreatedBy, &item.Content, &parentID, &likes, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return err
		}
		item.Likes = int(likes)
		_ = parentID
		items = append(items, item)
	}
	*comments = items
	return rows.Err()
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
