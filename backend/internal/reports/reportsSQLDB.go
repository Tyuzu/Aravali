// File: internal/reports/reportsSQLDB.go

package reports

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"scav/config"
	"scav/infra"
)

/* -------------------------
   Tables
------------------------- */

var (
	reportsTable       = config.Tables.ReportsTable
	appealsTable       = config.Tables.AppealsTable
	moderatorAppsTable = config.Tables.ModeratorApplications
)

func ensureReportsDB(app *infra.Deps) error {
	if app == nil || app.SQLDB == nil {
		return errors.New("database not initialized")
	}
	return nil
}

func buildReportColumns() string {
	return `reportid, userid, entity_type, entity_id, reason, status, created_at, updated_at, metadata`
}

func buildAppealColumns() string {
	return `appealid, reportid, userid, status, reason, created_at, updated_at, metadata`
}

/* -------------------------
   DB Wrappers
------------------------- */

func FindReportByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out *Report) error {
	if err := ensureReportsDB(app); err != nil {
		return err
	}
	if out == nil {
		return errors.New("nil result")
	}
	if query == "" {
		query = "1 = 1"
	}
	statement := fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT 1", buildReportColumns(), reportsTable, query)
	row := app.SQLDB.QueryRow(ctx, statement, args...)
	var reportID, userID, entityType, entityID, reason, status string
	var createdAt, updatedAt time.Time
	var metadata []byte
	if err := row.Scan(&reportID, &userID, &entityType, &entityID, &reason, &status, &createdAt, &updatedAt, &metadata); err != nil {
		return err
	}
	*out = Report{
		ReportID:   reportID,
		ReportedBy: userID,
		TargetType: entityType,
		TargetID:   entityID,
		Reason:     reason,
		Status:     status,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
	if len(metadata) > 0 {
		_ = metadata
	}
	return nil
}

func InsertReport(ctx context.Context, app *infra.Deps, report Report) error {
	if err := ensureReportsDB(app); err != nil {
		return err
	}
	statement := fmt.Sprintf(`INSERT INTO %s (reportid, userid, entity_type, entity_id, reason, status, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, reportsTable)
	_, err := app.SQLDB.Exec(ctx, statement,
		report.ReportID,
		report.ReportedBy,
		report.TargetType,
		report.TargetID,
		report.Reason,
		report.Status,
		time.Now().UTC(),
		time.Now().UTC(),
		map[string]any{},
	)
	return err
}

func FindReports(ctx context.Context, app *infra.Deps, query string, args []any, out *[]Report) error {
	if err := ensureReportsDB(app); err != nil {
		return err
	}
	if out == nil {
		return errors.New("nil result")
	}
	if query == "" {
		query = "1 = 1"
	}
	statement := fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY created_at DESC", buildReportColumns(), reportsTable, query)
	rows, err := app.SQLDB.Query(ctx, statement, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var reports []Report
	for rows.Next() {
		var reportID, userID, entityType, entityID, reason, status string
		var createdAt, updatedAt time.Time
		var metadata []byte
		if err := rows.Scan(&reportID, &userID, &entityType, &entityID, &reason, &status, &createdAt, &updatedAt, &metadata); err != nil {
			return err
		}
		reports = append(reports, Report{
			ReportID:   reportID,
			ReportedBy: userID,
			TargetType: entityType,
			TargetID:   entityID,
			Reason:     reason,
			Status:     status,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		})
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*out = reports
	return nil
}

func UpdateReportByID(ctx context.Context, app *infra.Deps, reportID string, update map[string]any) (int64, error) {
	if err := ensureReportsDB(app); err != nil {
		return 0, err
	}
	if len(update) == 0 {
		return 0, nil
	}
	parts := make([]string, 0, len(update))
	args := make([]any, 0, len(update)+1)
	for key, value := range update {
		parts = append(parts, fmt.Sprintf("%s = $%d", strings.TrimSpace(key), len(args)+1))
		args = append(args, value)
	}
	args = append(args, reportID)
	statement := fmt.Sprintf("UPDATE %s SET %s, updated_at = NOW() WHERE reportid = $%d", reportsTable, strings.Join(parts, ", "), len(args))
	tag, err := app.SQLDB.Exec(ctx, statement, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// Appeals
func FindAppealByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	if err := ensureReportsDB(app); err != nil {
		return err
	}
	if out == nil {
		return errors.New("nil result")
	}
	return nil
}

func InsertAppeal(ctx context.Context, app *infra.Deps, appeal any) error {
	if err := ensureReportsDB(app); err != nil {
		return err
	}
	return nil
}

func FindAppeals(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	if err := ensureReportsDB(app); err != nil {
		return err
	}
	return nil
}

func GetAppealByID(ctx context.Context, app *infra.Deps, appealID string, out any) error {
	if err := ensureReportsDB(app); err != nil {
		return err
	}
	return nil
}

func UpdateAppealByID(ctx context.Context, app *infra.Deps, appealID string, update map[string]any) (int64, error) {
	if err := ensureReportsDB(app); err != nil {
		return 0, err
	}
	return 0, nil
}

func setEntityDeletedFlagInDB(ctx context.Context, app *infra.Deps, table, idField, id string, deleted bool, by string) error {
	if err := ensureReportsDB(app); err != nil {
		return err
	}
	return nil
}
