// File: internal/notices/noticesSQLDB.go

package notices

import (
	"context"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
	"scav/utils"
)

var noticesTable = config.Tables.NoticesTable

func createNotice(ctx context.Context, app *infra.Deps, notice Notice) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if strings.TrimSpace(notice.NoticeID) == "" {
		notice.NoticeID = utils.GenerateRandomDigitString(13)
	}
	if notice.CreatedAt.IsZero() {
		notice.CreatedAt = time.Now()
	}
	if notice.UpdatedAt.IsZero() {
		notice.UpdatedAt = notice.CreatedAt
	}

	_, err := app.SQLDB.Exec(
		ctx,
		`INSERT INTO `+noticesTable+` (noticeid, title, content, userid, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		notice.NoticeID,
		notice.Title,
		notice.Content,
		notice.CreatedBy,
		notice.CreatedAt,
		notice.UpdatedAt,
	)
	return err
}

func findNoticeByID(ctx context.Context, app *infra.Deps, noticeID string) (Notice, error) {
	if app == nil || app.SQLDB == nil {
		return Notice{}, fmt.Errorf("notice %s not found", noticeID)
	}
	var n Notice
	var createdAt, updatedAt time.Time
	row := app.SQLDB.QueryRow(
		ctx,
		`SELECT noticeid, title, content, userid, created_at, updated_at FROM `+noticesTable+` WHERE noticeid = $1 LIMIT 1`,
		noticeID,
	)
	if err := row.Scan(&n.NoticeID, &n.Title, &n.Content, &n.CreatedBy, &createdAt, &updatedAt); err != nil {
		return Notice{}, err
	}
	n.CreatedAt = createdAt
	n.UpdatedAt = updatedAt
	return n, nil
}

func listNoticesWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any, out *[]Notice) error {
	if app == nil || app.SQLDB == nil || out == nil {
		return nil
	}
	if query == "" {
		query = "1 = 1"
	}

	sortBy := "created_at DESC"
	if v, ok := opts["sortBy"].(string); ok && strings.TrimSpace(v) != "" {
		sortBy = v
	}
	limit := 20
	if v, ok := opts["limit"].(int); ok && v > 0 {
		limit = v
	}
	if v, ok := opts["limit"].(int64); ok && v > 0 {
		limit = int(v)
	}
	page := 1
	if v, ok := opts["page"].(int); ok && v > 0 {
		page = v
	}
	if v, ok := opts["page"].(int64); ok && v > 0 {
		page = int(v)
	}
	offset := (page - 1) * limit

	rows, err := app.SQLDB.Query(
		ctx,
		`SELECT noticeid, title, content, userid, created_at, updated_at FROM `+noticesTable+` WHERE `+query+` ORDER BY `+sortBy+` LIMIT $1 OFFSET $2`,
		append(args, limit, offset)...,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	items := make([]Notice, 0)
	for rows.Next() {
		var n Notice
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&n.NoticeID, &n.Title, &n.Content, &n.CreatedBy, &createdAt, &updatedAt); err != nil {
			return err
		}
		n.CreatedAt = createdAt
		n.UpdatedAt = updatedAt
		items = append(items, n)
	}
	*out = items
	return rows.Err()
}

func buildNoticeQuery(entityType, entityID string) (string, []any) {
	if strings.TrimSpace(entityType) == "" || strings.TrimSpace(entityID) == "" {
		return "1 = 0", nil
	}

	return "entityType = $1 AND entityId = $2", []any{entityType, entityID}
}

func buildNoticeListOptions(page, limit int, sortBy string) map[string]any {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
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

func getNoticesPage(ctx context.Context, app *infra.Deps, entityType, entityID string, page, limit int, sortBy string) ([]Notice, error) {
	query, args := buildNoticeQuery(entityType, entityID)
	opts := buildNoticeListOptions(page, limit, sortBy)
	var notices []Notice
	if err := listNoticesWithOptions(ctx, app, query, args, opts, &notices); err != nil {
		return nil, err
	}
	if notices == nil {
		return []Notice{}, nil
	}
	return notices, nil
}

func buildNoticeSummary(notices []Notice) []map[string]any {
	if notices == nil {
		return []map[string]any{}
	}

	resp := make([]map[string]any, len(notices))
	for i, n := range notices {
		resp[i] = map[string]any{
			"noticeid":  n.NoticeID,
			"title":     n.Title,
			"summary":   n.Summary,
			"createdBy": n.CreatedBy,
			"createdAt": n.CreatedAt,
		}
	}
	return resp
}

func updateNoticeByID(ctx context.Context, app *infra.Deps, noticeID string, update map[string]any) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if strings.TrimSpace(noticeID) == "" {
		return fmt.Errorf("notice id is required")
	}
	if len(update) == 0 {
		return nil
	}

	setParts := make([]string, 0, len(update))
	args := make([]any, 0, len(update)+1)
	for key, value := range update {
		if key == "" || value == nil {
			continue
		}
		setParts = append(setParts, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value)
	}
	if len(setParts) == 0 {
		return nil
	}
	args = append(args, noticeID)
	_, err := app.SQLDB.Exec(ctx, `UPDATE `+noticesTable+` SET `+strings.Join(setParts, ", ")+` WHERE noticeid = $`+fmt.Sprintf("%d", len(args)), args...)
	return err
}

func deleteNoticeByID(ctx context.Context, app *infra.Deps, noticeID string) error {
	if app == nil || app.SQLDB == nil {
		return nil
	}
	if strings.TrimSpace(noticeID) == "" {
		return fmt.Errorf("notice id is required")
	}
	_, err := app.SQLDB.Exec(ctx, `DELETE FROM `+noticesTable+` WHERE noticeid = $1`, noticeID)
	return err
}
