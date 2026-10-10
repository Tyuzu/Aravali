// File: internal/verticals/notices/noticesSQLDB.go

package notices

import (
	"context"
	"strings"

	"scav/config"
	"scav/infra"
	"scav/infra/sqldb"
)

var noticesTable = config.Tables.NoticesTable

func createNotice(ctx context.Context, app *infra.Deps, notice Notice) error {
	return app.SQLDB.Insert(ctx, noticesTable, notice)
}

func findNoticeByID(ctx context.Context, app *infra.Deps, noticeID string) (Notice, error) {
	var notice Notice
	query := "noticeid = $1"
	args := []any{noticeID}
	err := app.SQLDB.FindOne(ctx, noticesTable, query, args, &notice)
	return notice, err
}

func listNoticesWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts sqldb.FindManyOptions, out *[]Notice) error {
	return app.SQLDB.FindManyWithOptions(ctx, noticesTable, query, args, opts, out)
}

func buildNoticeQuery(entityType, entityID string) (string, []any) {
	if strings.TrimSpace(entityType) == "" || strings.TrimSpace(entityID) == "" {
		return "1 = 0", nil
	}

	return "entityType = $1 AND entityId = $2", []any{entityType, entityID}
}

func buildNoticeListOptions(page, limit int, sortBy string) sqldb.FindManyOptions {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	orderBy := "createdAt DESC"
	if sortBy == "old" {
		orderBy = "createdAt ASC"
	}

	return sqldb.FindManyOptions{
		Limit:   int64(limit),
		Offset:  int64((page - 1) * limit),
		OrderBy: orderBy,
	}
}

func getNoticesPage(ctx context.Context, app *infra.Deps, entityType, entityID string, page, limit int, sortBy string) ([]Notice, error) {
	query, args := buildNoticeQuery(entityType, entityID)
	opts := buildNoticeListOptions(page, limit, sortBy)
	var notices []Notice
	if err := listNoticesWithOptions(ctx, app, query, args, opts, &notices); err != nil {
		return nil, err
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
	query := "noticeid = $1"
	args := []any{noticeID}
	_, err := app.SQLDB.Update(ctx, noticesTable, query, args, update)
	return err
}

func deleteNoticeByID(ctx context.Context, app *infra.Deps, noticeID string) error {
	query := "noticeid = $1"
	args := []any{noticeID}
	_, err := app.SQLDB.Delete(ctx, noticesTable, query, args)
	return err
}
