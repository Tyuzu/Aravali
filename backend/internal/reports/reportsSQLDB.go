package reports

import (
	"context"
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

func SQLFindReportByFilter(ctx context.Context, app *infra.Deps, filter map[string]any, out *Report) error {
	return app.DB.FindOne(ctx, reportsTable, filter, out)
}

func SQLInsertReport(ctx context.Context, app *infra.Deps, report Report) error {
	return app.DB.Insert(ctx, reportsTable, report)
}

func SQLFindReports(ctx context.Context, app *infra.Deps, filter map[string]any, out *[]Report) error {
	return app.DB.FindMany(ctx, reportsTable, filter, out)
}

func SQLUpdateReportByID(ctx context.Context, app *infra.Deps, reportID string, update map[string]any) (any, error) {
	return app.DB.Update(ctx, reportsTable, map[string]any{"reportid": reportID}, update)
}

// Appeals
func SQLFindAppealByFilter(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindOne(ctx, appealsTable, filter, out)
}

func SQLInsertAppeal(ctx context.Context, app *infra.Deps, appeal any) error {
	return app.DB.Insert(ctx, appealsTable, appeal)
}

func SQLFindAppeals(ctx context.Context, app *infra.Deps, filter map[string]any, out any) error {
	return app.DB.FindMany(ctx, appealsTable, filter, out)
}

func SQLGetAppealByID(ctx context.Context, app *infra.Deps, appealID string, out any) error {
	return app.DB.FindOne(ctx, appealsTable, map[string]any{"appealid": appealID}, out)
}

func SQLUpdateAppealByID(ctx context.Context, app *infra.Deps, appealID string, update map[string]any) (any, error) {
	return app.DB.Update(ctx, appealsTable, map[string]any{"appealid": appealID}, update)
}

func SQLsetEntityDeletedFlagInDB(ctx context.Context, app *infra.Deps, table, idField, id string, deleted bool, by string) error {
	now := time.Now().UTC()
	deletedAtVal := interface{}("")
	if deleted {
		deletedAtVal = now
	}

	_, err := app.DB.Update(
		ctx,
		table,
		map[string]any{idField: id},
		map[string]any{
			"deleted":   deleted,
			"deletedBy": by,
			"deletedAt": deletedAtVal,
		},
	)
	return err
}
