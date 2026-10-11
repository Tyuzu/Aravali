// File: internal/verticals/notices/noticesSQLDB.go

package notices

import (
	"context"
	"strings"

	"scav/config"
	"scav/infra"
)

var noticesTable = config.Tables.NoticesTable

func createNotice(ctx context.Context, app *infra.Deps, notice Notice) error {

}

func findNoticeByID(ctx context.Context, app *infra.Deps, noticeID string) (Notice, error) {
	var notice Notice
	return notice, err
}

func listNoticesWithOptions(ctx context.Context, app *infra.Deps, query string, args []any, opts map[string]any, out *[]Notice) error {

}

func buildNoticeQuery(entityType, entityID string) (string, []any) {
	if strings.TrimSpace(entityType) == "" || strings.TrimSpace(entityID) == "" {
		return "1 = 0", nil
	}

	return "entityType = $1 AND entityId = $2", []any{entityType, entityID}
}

func buildNoticeListOptions(page, limit int, sortBy string) map[string]any {

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

}

func deleteNoticeByID(ctx context.Context, app *infra.Deps, noticeID string) error {

}
