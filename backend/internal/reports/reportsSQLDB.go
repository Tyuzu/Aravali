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

func SQLFindReportByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out *Report) error {
	return app.SQLDB.FindOne(ctx, reportsTable, query, args, out)
}

func SQLInsertReport(ctx context.Context, app *infra.Deps, report Report) error {
	return app.SQLDB.Insert(ctx, reportsTable, report)
}

func SQLFindReports(ctx context.Context, app *infra.Deps, query string, args []any, out *[]Report) error {
	return app.SQLDB.FindMany(ctx, reportsTable, query, args, out)
}

func SQLUpdateReportByID(ctx context.Context, app *infra.Deps, reportID string, update map[string]any) (int64, error) {
	query := "reportid = $1"
	args := []any{reportID}

	return app.SQLDB.Update(ctx, reportsTable, query, args, update)
}

// Appeals
func SQLFindAppealByFilter(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindOne(ctx, appealsTable, query, args, out)
}

func SQLInsertAppeal(ctx context.Context, app *infra.Deps, appeal any) error {
	return app.SQLDB.Insert(ctx, appealsTable, appeal)
}

func SQLFindAppeals(ctx context.Context, app *infra.Deps, query string, args []any, out any) error {
	return app.SQLDB.FindMany(ctx, appealsTable, query, args, out)
}

func SQLGetAppealByID(ctx context.Context, app *infra.Deps, appealID string, out any) error {
	query := "appealid = $1"
	args := []any{appealID}

	return app.SQLDB.FindOne(ctx, appealsTable, query, args, out)
}

func SQLUpdateAppealByID(ctx context.Context, app *infra.Deps, appealID string, update map[string]any) (int64, error) {
	query := "appealid = $1"
	args := []any{appealID}

	return app.SQLDB.Update(ctx, appealsTable, query, args, update)
}

func SQLsetEntityDeletedFlagInDB(ctx context.Context, app *infra.Deps, table, idField, id string, deleted bool, by string) error {
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
