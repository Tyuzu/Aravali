// File: internal/reports/reportsSQLDB.go

package reports

import (
	"context"
	"fmt"
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

/* -------------------------
   DB Wrappers
------------------------- */

func FindReportByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out *Report) error {
	return app.SQLDB.FindOne(ctx, reportsTable, query, args, out)
}

func InsertReport(ctx context.Context, app *infra.Deps, report Report) error {
	return app.SQLDB.Insert(ctx, reportsTable, report)
}

func FindReports(ctx context.Context, app *infra.Deps, query string, args []any, out *[]Report) error {
	return app.SQLDB.FindMany(ctx, reportsTable, query, args, out)
}

func UpdateReportByID(ctx context.Context, app *infra.Deps, reportID string, update map[string]any) (int64, error) {
	query := "reportid = $1"
	args := []any{reportID}

	return app.SQLDB.Update(ctx, reportsTable, query, args, update)
}

// Appeals
func FindAppealByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindOne(ctx, appealsTable, query, args, out)
}

func InsertAppeal(ctx context.Context, app *infra.Deps, appeal any) error {
	return app.SQLDB.Insert(ctx, appealsTable, appeal)
}

func FindAppeals(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindMany(ctx, appealsTable, query, args, out)
}

func GetAppealByID(ctx context.Context, app *infra.Deps, appealID string, out any) error {
	query := "appealid = $1"
	args := []any{appealID}

	return app.SQLDB.FindOne(ctx, appealsTable, query, args, out)
}

func UpdateAppealByID(ctx context.Context, app *infra.Deps, appealID string, update map[string]any) (int64, error) {
	query := "appealid = $1"
	args := []any{appealID}

	return app.SQLDB.Update(ctx, appealsTable, query, args, update)
}

func setEntityDeletedFlagInDB(ctx context.Context, app *infra.Deps, table, idField, id string, deleted bool, by string) error {
	now := time.Now().UTC()
	var deletedAtVal any = nil
	if deleted {
		deletedAtVal = now
	}

	query := fmt.Sprintf("%s = $1", idField)
	args := []any{id}

	update := map[string]any{
		"deleted":   deleted,
		"deletedBy": by,
		"deletedAt": deletedAtVal,
	}

	_, err := app.SQLDB.Update(
		ctx,
		table,
		query,
		args,
		update,
	)
	return err
}
